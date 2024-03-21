// Package loop provides helpers for loop operation and iterator implementations
package loop

import (
	"errors"
	"unsafe"

	"github.com/m4gshm/gollections/break/predicate/always"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
)

// ErrBreak is the 'break' statement of the For, Track methods.
var ErrBreak = c.ErrBreak

// Of wrap the elements by loop function
func Of[T any](elements ...T) func() (e T, ok bool, err error) {
	l := len(elements)
	i := 0
	if l == 0 || i < 0 || i >= l {
		return func() (e T, ok bool, err error) { return e, false, nil }
	}
	return func() (e T, ok bool, err error) {
		if i < l {
			e, ok = elements[i], true
			i++
		}
		return e, ok, nil
	}
}

// New is the main breakable loop constructor
func New[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) Loop[T] {
	return func() (out T, ok bool, err error) {
		if ok := hasNext(source); !ok {
			return out, false, nil
		}
		out, err = getNext(source)
		return out, err == nil, err
	}
}

// From wrap the next loop to a breakable loop
func From[T any](next func() (T, bool)) Loop[T] {
	return func() (T, bool, error) {
		e, ok := next()
		return e, ok, nil
	}
}

// To transforms a breakable loop to a simple loop.
// The errConsumer is a function that is called when an error occurs.
func To[T any](next func() (T, bool, error), errConsumer func(error)) func() (T, bool) {
	return func() (T, bool) {
		e, ok, err := next()
		if err != nil {
			errConsumer(err)
			return e, false
		}
		return e, ok
	}
}

// For applies the 'walker' function for the elements retrieved by the 'next' function. Return the c.ErrBreak to stop
func For[T any](next func() (T, bool, error), walker func(T) error) error {
	for {
		if v, ok, err := next(); err != nil || !ok {
			return err
		} else if err := walker(v); err != nil {
			return brk(err)
		}
	}
}

// ForFiltered applies the 'walker' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func ForFiltered[T any](next func() (T, bool, error), walker func(T) error, predicate func(T) bool) error {
	for {
		if v, ok, err := next(); err != nil || !ok {
			return err
		} else if ok := predicate(v); ok {
			if err := walker(v); err != nil {
				return brk(err)
			}
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool, error), predicate func(T) bool) (T, bool, error) {
	for {
		if out, ok, err := next(); err != nil || !ok {
			return out, false, err
		} else if ok := predicate(out); ok {
			return out, true, nil
		}
	}
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any](next func() (T, bool, error), predicate func(T) (bool, error)) (T, bool, error) {
	for {
		if out, ok, err := next(); err != nil || !ok {
			return out, false, err
		} else if ok, err := predicate(out); err != nil || ok {
			return out, ok, err
		}
	}
}

// Track applies the 'tracker' function to position/element pairs retrieved by the 'next' function. Return the c.ErrBreak to stop tracking..
func Track[I, T any](next func() (I, T, bool, error), tracker func(I, T) error) error {
	for {
		if p, v, ok, err := next(); err != nil || !ok {
			return err
		} else if err := tracker(p, v); err != nil {
			return brk(err)
		}
	}
}

// Slice collects the elements retrieved by the 'next' function into a slice
func Slice[T any](next func() (T, bool, error)) (out []T, err error) {
	if next == nil {
		return nil, nil
	}
	for {
		v, ok, err := next()
		if ok {
			out = append(out, v)
		}
		if !ok || err != nil {
			return out, err
		}
	}
}

// SliceCap collects the elements retrieved by the 'next' function into a new slice with predefined capacity
func SliceCap[T any](next func() (T, bool, error), cap int) (out []T, err error) {
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(next, out)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func Append[T any, TS ~[]T](next func() (T, bool, error), out TS) (TS, error) {
	for v, ok, err := next(); ok; v, ok, err = next() {
		if err != nil {
			return out, err
		}
		out = append(out, v)
	}
	return out, nil
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function
func Reduce[T any](next func() (T, bool, error), merger func(T, T) T) (out T, e error) {
	v, ok, err := next()
	if err != nil || !ok {
		return out, err
	}
	out = v
	for {
		v, ok, err := next()
		if err != nil || !ok {
			return out, err
		}
		out = merger(out, v)
	}
}

// Reducee reduces the elements retrieved by the 'next' function into an one using the 'merge' function
func Reducee[T any](next func() (T, bool, error), merger func(T, T) (T, error)) (out T, e error) {
	v, ok, err := next()
	if err != nil || !ok {
		return out, err
	}
	out = v
	for {
		if v, ok, err := next(); err != nil || !ok {
			return out, err
		} else if out, err = merger(out, v); err != nil {
			return out, err
		}
	}
}

// Sum returns the sum of all elements
func Sum[T c.Summable](next func() (T, bool, error)) (T, error) {
	return Reduce(next, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func HasAny[T any](next func() (T, bool, error), predicate func(T) bool) (bool, error) {
	_, ok, err := First(next, predicate)
	return ok, err
}

// HasAnyy finds the first element that satisfies the 'predicate' function condition and returns true if successful
func HasAnyy[T any](next func() (T, bool, error), predicate func(T) (bool, error)) (bool, error) {
	_, ok, err := Firstt(next, predicate)
	return ok, err
}

// Contains  finds the first element that equal to the example and returns true
func Contains[T comparable](next func() (T, bool, error), example T) (bool, error) {
	for {
		if one, ok, err := next(); err != nil || !ok {
			return false, err
		} else if one == example {
			return true, nil
		}
	}
}

// Conv instantiates an iterator that converts elements with a converter and returns them.
func Conv[From, To any](next func() (From, bool, error), converter func(From) (To, error)) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool, err error) {
		v, ok, err := next()
		if err != nil || !ok {
			return t, ok, err
		}
		vc, err := converter(v)
		return vc, err == nil, err
	}
}

// Convert instantiates an iterator that converts elements with a converter and returns them.
func Convert[From, To any](next func() (From, bool, error), converter func(From) To) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool, err error) {
		if v, ok, err := next(); err != nil || !ok {
			return t, ok, err
		} else {
			return converter(v), true, nil
		}
	}
}

// ConvCheck is similar to ConvertFilt, but it checks and transforms elements together
func ConvCheck[From, To any](next func() (From, bool, error), converter func(from From) (To, bool, error)) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool, err error) {
		for {
			if v, ok, err := next(); err != nil || !ok {
				return t, false, err
			} else if vc, ok, err := converter(v); err != nil || ok {
				return vc, ok, err
			}
		}
	}
}

// ConvertCheck is similar to ConvFilt, but it checks and transforms elements together
func ConvertCheck[From, To any](next func() (From, bool, error), converter func(from From) (To, bool)) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool, err error) {
		for {
			if e, ok, err := next(); err != nil || !ok {
				return t, false, err
			} else if t, ok := converter(e); ok {
				return t, ok, err
			}
		}
	}
}

// FiltAndConv returns a stream that filters source elements and converts them
func FiltAndConv[From, To any](next func() (From, bool, error), filter func(From) (bool, error), converter func(From) (To, error)) Loop[To] {
	return FilterConvertFilter(next, filter, converter, always.True[To])
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any](next func() (From, bool, error), filter func(From) bool, converter func(From) To) Loop[To] {
	return FilterConvertFilter(next, func(f From) (bool, error) { return filter(f), nil }, func(f From) (To, error) { return converter(f), nil }, always.True[To])
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[From, To any](next func() (From, bool, error), filter func(From) (bool, error), converter func(From) (To, error), filterTo func(To) (bool, error)) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool, err error) {
		for {
			if f, ok, err := Firstt(next, filter); err != nil || !ok {
				return t, false, err
			} else if cf, err := converter(f); err != nil {
				return t, false, err
			} else if ok, err := filterTo(cf); err != nil || !ok {
				return t, false, err
			} else {
				return cf, true, nil
			}
		}
	}
}

// ConvertAndFilter additionally filters 'To' elements
func ConvertAndFilter[From, To any](next func() (From, bool, error), converter func(From) (To, error), filter func(To) (bool, error)) Loop[To] {
	return FilterConvertFilter(next, always.True[From], converter, filter)
}

// Flatt instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](next func() (From, bool, error), flattener func(From) ([]To, error)) Loop[To] {
	if next == nil {
		return nil
	}
	var (
		elemSizeTo      uintptr = notsafe.GetTypeSize[To]()
		arrayTo         unsafe.Pointer
		indexTo, sizeTo int
	)
	return func() (t To, ok bool, err error) {
		if sizeTo > 0 {
			if indexTo < sizeTo {
				indexTo++
				return *(*To)(notsafe.GetArrayElemRef(arrayTo, indexTo, elemSizeTo)), true, nil
			}
			indexTo = 0
			arrayTo = nil
			sizeTo = 0
		}
		for {
			if v, ok, err := next(); err != nil || !ok {
				return t, ok, err
			} else if elementsTo, err := flattener(v); err != nil {
				return t, false, err
			} else if len(elementsTo) > 0 {
				indexTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				arrayTo = unsafe.Pointer(header.Data)
				sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(arrayTo, 0, elemSizeTo)), true, nil
			}
		}
	}
	// return &FlatIter[From, To]{next: next, flattener: flattener, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Flat instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flat[From, To any](next func() (From, bool, error), flattener func(From) []To) Loop[To] {
	if next == nil {
		return nil
	}
	var (
		elemSizeTo      uintptr = notsafe.GetTypeSize[To]()
		arrayTo         unsafe.Pointer
		indexTo, sizeTo int
	)
	return func() (t To, ok bool, err error) {
		if sizeTo > 0 {
			if indexTo < sizeTo {
				i := indexTo
				indexTo++
				return *(*To)(notsafe.GetArrayElemRef(arrayTo, i, elemSizeTo)), true, nil
			}
			indexTo = 0
			arrayTo = nil
			sizeTo = 0
		}
		for {
			if v, ok, err := next(); err != nil {
				return t, false, err
			} else if !ok {
				return t, false, nil
			} else if elementsTo := flattener(v); len(elementsTo) > 0 {
				indexTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				arrayTo = unsafe.Pointer(header.Data)
				sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(arrayTo, 0, elemSizeTo)), true, nil
			}
		}
	}
}

// FiltAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FiltAndFlat[From, To any](next func() (From, bool, error), filter func(From) (bool, error), flattener func(From) ([]To, error)) Loop[To] {
	return FiltFlattFilt(next, filter, flattener, always.True[To])
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any](next func() (From, bool, error), filter func(From) bool, flattener func(From) []To) Loop[To] {
	return FiltFlattFilt(next, func(f From) (bool, error) { return filter(f), nil }, func(f From) ([]To, error) { return flattener(f), nil }, always.True[To])
}

// FlatAndFilt extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlatAndFilt[From, To any](next func() (From, bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) Loop[To] {
	return FiltFlattFilt(next, always.True[From], flattener, filterTo)
}

// FlattAndFilter extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlattAndFilter[From, To any](next func() (From, bool, error), flattener func(From) []To, filterTo func(To) bool) Loop[To] {
	return FiltFlattFilt(next, always.True[From], func(f From) ([]To, error) { return flattener(f), nil }, func(t To) (bool, error) { return filterTo(t), nil })
}

// FiltFlattFilt filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FiltFlattFilt[From, To any](next func() (From, bool, error), filterFrom func(From) (bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) Loop[To] {
	if next == nil {
		return nil
	}
	var (
		elemSizeTo      uintptr = notsafe.GetTypeSize[To]()
		arrayTo         unsafe.Pointer
		indexTo, sizeTo int
	)
	return func() (t To, ok bool, err error) {
		for {
			if sizeTo > 0 {
				if indexTo < sizeTo {
					i := indexTo
					indexTo++
					t = *(*To)(notsafe.GetArrayElemRef(arrayTo, i, elemSizeTo))
					if ok, err := filterTo(t); err != nil {
						return t, false, err
					} else if ok {
						return t, true, nil
					}
				}
				indexTo = 0
				arrayTo = nil
				sizeTo = 0
			}

			if v, ok, err := next(); err != nil || !ok {
				return t, false, err
			} else if ok, err := filterFrom(v); err != nil {
				return t, false, err
			} else if ok {
				if elementsTo, err := flattener(v); err != nil {
					return t, false, err
				} else if len(elementsTo) > 0 {
					indexTo = 1
					header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
					arrayTo = unsafe.Pointer(header.Data)
					sizeTo = header.Len
					t = *(*To)(notsafe.GetArrayElemRef(arrayTo, 0, elemSizeTo))
					if ok, err := filterTo(t); err != nil || ok {
						return t, ok, err
					}
				}
			}
		}
	}
	// return &FlattFiltIter[From, To]{next: next, filterFrom: filterFrom, flattener: flattener, filterTo: filterTo, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FilterFlatFilter filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FilterFlatFilter[From, To any](next func() (From, bool, error), filterFrom func(From) bool, flattener func(From) []To, filterTo func(To) bool) Loop[To] {
	if next == nil {
		return nil
	}
	var (
		elemSizeTo      uintptr = notsafe.GetTypeSize[To]()
		arrayTo         unsafe.Pointer
		indexTo, sizeTo int
	)
	return func() (t To, ok bool, err error) {
		for {
			if sizeTo > 0 {
				if indexTo < sizeTo {
					i := indexTo
					indexTo++
					tv := *(*To)(notsafe.GetArrayElemRef(arrayTo, i, elemSizeTo))
					if ok := filterTo(tv); ok {
						return tv, true, nil
					}
				}
				indexTo = 0
				arrayTo = nil
				sizeTo = 0
			}

			if fv, ok, err := next(); err != nil || !ok {
				return t, false, err
			} else if ok := filterFrom(fv); ok {
				if elementsTo := flattener(fv); len(elementsTo) > 0 {
					indexTo = 1
					header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
					arrayTo = unsafe.Pointer(header.Data)
					sizeTo = header.Len
					tv := *(*To)(notsafe.GetArrayElemRef(arrayTo, 0, elemSizeTo))
					if ok := filterTo(tv); ok {
						return tv, true, nil
					}
				}
			}
		}
	}
}

// Filt creates an iterator that checks elements by the 'filter' function and returns successful ones.
func Filt[T any](next func() (T, bool, error), filter func(T) (bool, error)) Loop[T] {
	if next == nil {
		return nil
	}
	return func() (T, bool, error) {
		return Firstt(next, filter)
	}
}

// Filter creates an iterator that checks elements by the 'filter' function and returns successful ones.
func Filter[T any](next func() (T, bool, error), filter func(T) bool) Loop[T] {
	if next == nil {
		return nil
	}
	return func() (T, bool, error) {
		return First(next, filter)
	}
}

// NotNil creates an iterator that filters nullable elements.
func NotNil[T any](next func() (*T, bool, error)) Loop[*T] {
	return Filt(next, as.ErrTail(not.Nil[T]))
}

// PtrVal creates an iterator that transform pointers to the values referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any](next func() (*T, bool, error)) Loop[T] {
	return Convert(next, convert.PtrVal[T])
}

// NoNilPtrVal creates an iterator that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any](next func() (*T, bool, error)) Loop[T] {
	return ConvertCheck(next, convert.NoNilPtrVal[T])
}

// KeyValue transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValue[T any, K, V any](next func() (T, bool, error), keyExtractor func(T) K, valExtractor func(T) V) func() (K, V, bool, error) {
	return KeyValuee(next, as.ErrTail(keyExtractor), as.ErrTail(valExtractor))
}

// KeyValuee transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValuee[T any, K, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) func() (K, V, bool, error) {
	if next == nil {
		return nil
	}
	return func() (key K, value V, ok bool, err error) {
		if elem, nextOk, err := next(); err != nil || !nextOk {
			return key, value, false, err
		} else if key, err = keyExtractor(elem); err == nil {
			value, err = valExtractor(elem)
			return key, value, err == nil, err
		}
		return key, value, false, nil
	}
}

// KeysValues transforms iterable elements to key/value iterator based on applying multiple keys, values extractor to the elements
func KeysValues[T, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valsExtractor func(T) ([]V, error)) func() (K, V, bool, error) {
	if next == nil {
		return nil
	}
	var (
		keys   []K
		values []V
		ki, vi int
	)
	return func() (key K, value V, ok bool, err error) {
		for !ok {
			var (
				keysLen, valuesLen         = len(keys), len(values)
				lastKeyIndex, lastValIndex = keysLen - 1, valuesLen - 1
			)
			if keysLen > 0 && ki >= 0 && ki <= lastKeyIndex {
				key = keys[ki]
				ok = true
			}
			if valuesLen > 0 && vi >= 0 && vi <= lastValIndex {
				value = values[vi]
				ok = true
			}
			if ok {
				if ki < lastKeyIndex {
					ki++
				} else if vi < lastValIndex {
					ki = 0
					vi++
				} else {
					keys, values = nil, nil
				}
			} else if elem, nextOk, err := next(); err != nil {
				return key, value, ok, err
			} else if nextOk {
				keys, err = keysExtractor(elem)
				if err == nil {
					values, err = valsExtractor(elem)
				}
				if err != nil {
					break
				}
				ki, vi = 0, 0
			} else {
				keys, values = nil, nil
				break
			}
		}
		return key, value, ok, nil
	}
	// return NewMultipleKeyValuer(next, keysExtractor, valsExtractor)
}

// KeysValue transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func KeysValue[T, K, V any](next func() (T, bool, error), keysExtractor func(T) []K, valExtractor func(T) V) func() (K, V, bool, error) {
	return KeysValues(next, as.ErrTail(keysExtractor), convSlice(as.ErrTail(valExtractor)))
}

// KeysValuee transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func KeysValuee[T, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valExtractor func(T) (V, error)) func() (K, V, bool, error) {
	return KeysValues(next, keysExtractor, convSlice(valExtractor))
}

// KeyValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func KeyValues[T, K, V any](next func() (T, bool, error), keyExtractor func(T) K, valsExtractor func(T) []V) func() (K, V, bool, error) {
	return KeysValues(next, convSlice(as.ErrTail(keyExtractor)), as.ErrTail(valsExtractor))
}

// KeyValuess transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func KeyValuess[T, K, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valsExtractor func(T) ([]V, error)) func() (K, V, bool, error) {
	return KeysValues(next, convSlice(keyExtractor), valsExtractor)
}

// ExtraVals transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraVals[T, V any](next func() (T, bool, error), valsExtractor func(T) []V) func() (T, V, bool, error) {
	return KeyValues(next, as.Is[T], valsExtractor)
}

// ExtraValss transforms iterable elements to key/value iterator based on applying values extractor to the elements
func ExtraValss[T, V any](next func() (T, bool, error), valsExtractor func(T) ([]V, error)) func() (T, V, bool, error) {
	return KeyValuess(next, as.ErrTail(as.Is[T]), valsExtractor)
}

// ExtraKeys transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeys[T, K any](next func() (T, bool, error), keysExtractor func(T) []K) func() (K, T, bool, error) {
	return KeysValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyss transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeyss[T, K any](next func() (T, bool, error), keyExtractor func(T) (K, error)) func() (K, T, bool, error) {
	return KeyValuess(next, keyExtractor, as.ErrTail(convert.AsSlice[T]))
}

// ExtraKey transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKey[T, K any](next func() (T, bool, error), keysExtractor func(T) K) func() (K, T, bool, error) {
	return KeyValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyy transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeyy[T, K any](next func() (T, bool, error), keyExtractor func(T) (K, error)) func() (K, T, bool, error) {
	return KeyValuee[T, K](next, keyExtractor, as.ErrTail(as.Is[T]))
}

// ExtraValue transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraValue[T, V any](next func() (T, bool, error), valueExtractor func(T) V) func() (T, V, bool, error) {
	return KeyValue(next, as.Is[T], valueExtractor)
}

// ExtraValuee transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraValuee[T, V any](next func() (T, bool, error), valExtractor func(T) (V, error)) func() (T, V, bool, error) {
	return KeyValuee[T, T, V](next, as.ErrTail(as.Is[T]), valExtractor)
}

// Group converts elements retrieved by the 'next' function into a map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to an value.
func Group[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) K, valExtractor func(T) V) (map[K][]V, error) {
	return Groupp(next, as.ErrTail(keyExtractor), as.ErrTail(valExtractor))
}

// Groupp converts elements retrieved by the 'next' function into a map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to an value.
func Groupp[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K][]V, error) {
	return ToMapResolvv(next, keyExtractor, valExtractor, func(ok bool, k K, rv []V, v V) ([]V, error) {
		return resolv.Slice(ok, k, rv, v), nil
	})
}

// GroupByMultiple converts elements retrieved by the 'next' function into a map, extracting multiple keys, values per each element applying the 'keysExtractor' and 'valsExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valsExtractor retrieves one or more values per element.
func GroupByMultiple[T any, K comparable, V any](next func() (T, bool, error), keysExtractor func(T) []K, valsExtractor func(T) []V) (map[K][]V, error) {
	groups := map[K][]V{}
	for {
		if e, ok, err := next(); err != nil || !ok {
			return groups, err
		} else if keys, vals := keysExtractor(e), valsExtractor(e); len(keys) == 0 {
			var key K
			for _, v := range vals {
				initGroup(key, v, groups)
			}
		} else {
			for _, key := range keys {
				if len(vals) == 0 {
					var v V
					initGroup(key, v, groups)
				} else {
					for _, v := range vals {
						initGroup(key, v, groups)
					}
				}
			}
		}
	}
}

// GroupByMultipleKeys converts elements retrieved by the 'next' function into a map, extracting multiple keys, one value per each element applying the 'keysExtractor' and 'valExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valExtractor converts an element to a value.
func GroupByMultipleKeys[T any, K comparable, V any](next func() (T, bool, error), keysExtractor func(T) []K, valExtractor func(T) V) (map[K][]V, error) {
	groups := map[K][]V{}
	for {
		if e, ok, err := next(); err != nil || !ok {
			return groups, err
		} else if keys, v := keysExtractor(e), valExtractor(e); len(keys) == 0 {
			var key K
			initGroup(key, v, groups)
		} else {
			for _, key := range keys {
				initGroup(key, v, groups)
			}
		}
	}
}

// GroupByMultipleValues converts elements retrieved by the 'next' function into a map, extracting one key, multiple values per each element applying the 'keyExtractor' and 'valsExtractor' functions.
// The keyExtractor converts an element to a key.
// The valsExtractor retrieves one or more values per element.
func GroupByMultipleValues[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) K, valsExtractor func(T) []V) (map[K][]V, error) {
	groups := map[K][]V{}
	for {
		if e, ok, err := next(); err != nil || !ok {
			return groups, err
		} else if key, vals := keyExtractor(e), valsExtractor(e); len(vals) == 0 {
			var v V
			initGroup(key, v, groups)
		} else {
			for _, v := range vals {
				initGroup(key, v, groups)
			}
		}
	}
}

func initGroup[T any, K comparable, TS ~[]T](key K, e T, groups map[K]TS) {
	groups[key] = append(groups[key], e)
}

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) (map[K]V, error) {
	return ToMapp(From(next), as.ErrTail(keyExtractor), as.ErrTail(valExtractor))
}

// ToMapp collects key\value elements to a map by iterating over the elements
func ToMapp[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K]V, error) {
	return ToMapResolvv(next, keyExtractor, valExtractor, func(ok bool, k K, rv V, v V) (V, error) { return resolv.First(ok, k, rv, v), nil })
}

// ToMapResolvv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolvv[T any, K comparable, V, VR any](
	next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error),
	resolver func(bool, K, VR, V) (VR, error),
) (m map[K]VR, err error) {
	m = map[K]VR{}
	for {
		if e, ok, err := next(); err != nil || !ok {
			return m, err
		} else if k, err := keyExtractor(e); err != nil {
			return m, err
		} else if v, err := valExtractor(e); err != nil {
			return m, err
		} else {
			exists, ok := m[k]
			if m[k], err = resolver(ok, k, exists, v); err != nil {
				return m, err
			}
		}
	}
}

// ConvertAndReduce converts each elements and merges them into one
func ConvertAndReduce[From, To any](next func() (From, bool, error), converter func(From) To, merger func(To, To) To) (out To, err error) {
	return Reduce(Convert(next, converter), merger)
}

// ConvAndReduce converts each elements and merges them into one
func ConvAndReduce[From, To any](next func() (From, bool, error), converter func(From) (To, error), merger func(To, To) To) (out To, err error) {
	return Reduce(Conv(next, converter), merger)
}

func brk(err error) error {
	if errors.Is(err, c.ErrBreak) {
		return nil
	}
	return err
}

func convSlice[T, V any](conv func(T) (V, error)) func(t T) ([]V, error) {
	return func(t T) ([]V, error) {
		v, err := conv(t)
		if err != nil {
			return nil, err
		}
		return convert.AsSlice(v), nil
	}
}
