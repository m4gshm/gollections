// Package loop provides helpers for loop operation and iterator implementations
package loop

import (
	"errors"
	"unsafe"

	breakkvloop "github.com/m4gshm/gollections/break/kv/loop"
	"github.com/m4gshm/gollections/break/predicate/always"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
)

// Break is the 'break' statement of the For, Track methods.
var Break = c.Break

// S wrap the elements by loop function.
func S[TS ~[]T, T any](elements TS) Loop[T] {
	return Of(elements...)
}

// Of wrap the elements by loop function.
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
	if next == nil {
		return nil
	}
	return func() (T, bool, error) {
		e, ok := next()
		return e, ok, nil
	}
}

// To transforms a breakable loop to a simple loop.
// The errConsumer is a function that is called when an error occurs.
func To[T any](next func() (T, bool, error), errConsumer func(error)) func() (T, bool) {
	if next == nil {
		return nil
	}
	return func() (T, bool) {
		e, ok, err := next()
		if err != nil {
			errConsumer(err)
			return e, false
		}
		return e, ok
	}
}

// All is an adapter for the next function for iterating by `for ... range`.
func All[T any](next func() (T, bool, error), consumer func(T, error) bool) {
	if next == nil {
		return
	}
	for {
		v, ok, err := next()
		if !ok {
			if err != nil {
				consumer(v, err)
			}
			break
		} else {
			consumer(v, err)
		}
	}
}

// For applies the 'consumer' function for the elements retrieved by the 'next' function until the consumer returns the c.Break to stop.
func For[T any](next func() (T, bool, error), consumer func(T) error) error {
	if next == nil {
		return nil
	}
	for {
		if v, ok, err := next(); err != nil || !ok {
			return err
		} else if err := consumer(v); err != nil {
			return brk(err)
		}
	}
}

// ForFiltered applies the 'consumer' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func ForFiltered[T any](next func() (T, bool, error), consumer func(T) error, predicate func(T) bool) error {
	if next == nil {
		return nil
	}
	for {
		if v, ok, err := next(); err != nil || !ok {
			return err
		} else if ok := predicate(v); ok {
			if err := consumer(v); err != nil {
				return brk(err)
			}
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool, error), predicate func(T) bool) (out T, ok bool, err error) {
	if next == nil {
		return out, false, nil
	}
	for {
		if out, ok, err = next(); err != nil || !ok {
			return out, false, err
		} else if ok := predicate(out); ok {
			return out, true, nil
		}
	}
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any](next func() (T, bool, error), predicate func(T) (bool, error)) (out T, ok bool, err error) {
	if next == nil {
		return out, false, nil
	}
	for {
		if out, ok, err := next(); err != nil || !ok {
			return out, false, err
		} else if ok, err := predicate(out); err != nil || ok {
			return out, ok, err
		}
	}
}

// Track applies the 'consumer' function to position/element pairs retrieved by the 'next' function until the consumer returns the c.Break to stop.
func Track[I, T any](next func() (I, T, bool, error), consumer func(I, T) error) error {
	return breakkvloop.Track(next, consumer)
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
	if next == nil {
		return nil, nil
	}
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(next, out)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func Append[T any, TS ~[]T](next func() (T, bool, error), out TS) (TS, error) {
	if next == nil {
		return out, nil
	}
	for v, ok, err := next(); ok; v, ok, err = next() {
		if err != nil {
			return out, err
		}
		out = append(out, v)
	}
	return out, nil
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func Reduce[T any](next func() (T, bool, error), merge func(T, T) T) (T, error) {
	result, _, err := ReduceOK(next, merge)
	return result, err
}

// ReduceOK reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceOK[T any](next func() (T, bool, error), merge func(T, T) T) (result T, ok bool, err error) {
	if next == nil {
		return result, false, nil
	}
	if result, ok, err = next(); err != nil || !ok {
		return result, ok, err
	}
	result, err = Accum(result, next, merge)
	return result, true, err
}

// Reducee reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func Reducee[T any](next func() (T, bool, error), merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(next, merge)
	return result, err
}

// ReduceeOK reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceeOK[T any](next func() (T, bool, error), merge func(T, T) (T, error)) (result T, ok bool, err error) {
	if next == nil {
		return result, false, nil
	}
	if result, ok, err = next(); err != nil || !ok {
		return result, ok, err
	}
	result, err = Accumm(result, next, merge)
	return result, true, err
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func Accum[T any](first T, next func() (T, bool, error), merge func(T, T) T) (accumulator T, err error) {
	accumulator = first
	if next == nil {
		return accumulator, nil
	}
	for {
		v, ok, err := next()
		if err != nil {
			return accumulator, err
		} else if !ok {
			return accumulator, nil
		}
		accumulator = merge(accumulator, v)
	}
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func Accumm[T any](first T, next func() (T, bool, error), merge func(T, T) (T, error)) (accumulator T, err error) {
	accumulator = first
	if next == nil {
		return accumulator, nil
	}
	for {
		if v, ok, err := next(); err != nil || !ok {
			return accumulator, err
		} else if accumulator, err = merge(accumulator, v); err != nil {
			return accumulator, err
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
	if next == nil {
		return false, nil
	}
	for {
		if one, ok, err := next(); err != nil || !ok {
			return false, err
		} else if one == example {
			return true, nil
		}
	}
}

// Conv creates a loop that applies the 'converter' function to iterable elements.
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

// Convert creates a loop that applies the 'converter' function to iterable elements.
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

// ConvOK creates a loop that applies the 'converter' function to iterable elements.
// The converter may returns converted value or ok=false to exclude the value from the loop.
// It may also return an error to abort the loop.
func ConvOK[From, To any](next func() (From, bool, error), converter func(from From) (to To, ok bool, err error)) Loop[To] {
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

// ConvertOK creates a loop that applies the 'converter' function to iterable elements.
// The converter may returns a value or ok=false to exclude the value from the loop.
func ConvertOK[From, To any](next func() (From, bool, error), converter func(from From) (To, bool)) Loop[To] {
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

// FiltAndConv creates a loop that filters source elements and converts them
func FiltAndConv[From, To any](next func() (From, bool, error), filter func(From) (bool, error), converter func(From) (To, error)) Loop[To] {
	return FilterConvertFilter(next, filter, converter, always.True[To])
}

// FilterAndConvert creates a loop that filters source elements and converts them
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

// Flatt converts a two-dimensional loop in an one-dimensional one.
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
}

// Flat converts a two-dimensional loop in an one-dimensional one.
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

// Filt creates a loop that checks elements by the 'filter' function and returns successful ones.
func Filt[T any](next func() (T, bool, error), filter func(T) (bool, error)) Loop[T] {
	if next == nil {
		return nil
	}
	return func() (T, bool, error) {
		return Firstt(next, filter)
	}
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func Filter[T any](next func() (T, bool, error), filter func(T) bool) Loop[T] {
	if next == nil {
		return nil
	}
	return func() (T, bool, error) {
		return First(next, filter)
	}
}

// NotNil creates a loop that filters nullable elements.
func NotNil[T any](next func() (*T, bool, error)) Loop[*T] {
	return Filt(next, as.ErrTail(not.Nil[T]))
}

// PtrVal creates a loop that transform pointers to the values referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any](next func() (*T, bool, error)) Loop[T] {
	return Convert(next, convert.PtrVal[T])
}

// NoNilPtrVal creates a loop that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any](next func() (*T, bool, error)) Loop[T] {
	return ConvertOK(next, convert.NoNilPtrVal[T])
}

// KeyValue transforms a loop to the key/value loop based on applying key, value extractors to the elements
func KeyValue[T any, K, V any](next func() (T, bool, error), keyExtractor func(T) K, valExtractor func(T) V) breakkvloop.Loop[K, V] {
	return KeyValuee(next, as.ErrTail(keyExtractor), as.ErrTail(valExtractor))
}

// KeyValuee transforms a loop to the key/value loop based on applying key, value extractors to the elements
func KeyValuee[T any, K, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) breakkvloop.Loop[K, V] {
	if next == nil {
		return nil
	}
	return func() (key K, value V, ok bool, err error) {
		if elem, nextOk, err := next(); err != nil || !nextOk {
			return key, value, false, err
		} else if key, err = keyExtractor(elem); err == nil {
			value, err = valExtractor(elem)
			return key, value, true, err
		}
		return key, value, false, nil
	}
}

// KeysValues transforms a loop to the key/value loop based on applying multiple keys, values extractor to the elements
func KeysValues[T, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valsExtractor func(T) ([]V, error)) breakkvloop.Loop[K, V] {
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

// KeysValue transforms a loop to the key/value loop based on applying key, value extractor to the elements
func KeysValue[T, K, V any](next func() (T, bool, error), keysExtractor func(T) []K, valExtractor func(T) V) breakkvloop.Loop[K, V] {
	return KeysValues(next, as.ErrTail(keysExtractor), convSlice(as.ErrTail(valExtractor)))
}

// KeysValuee transforms a loop to the key/value loop based on applying key, value extractor to the elements
func KeysValuee[T, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valExtractor func(T) (V, error)) breakkvloop.Loop[K, V] {
	return KeysValues(next, keysExtractor, convSlice(valExtractor))
}

// KeyValues transforms a loop to the key/value loop based on applying key, value extractor to the elements
func KeyValues[T, K, V any](next func() (T, bool, error), keyExtractor func(T) K, valsExtractor func(T) []V) breakkvloop.Loop[K, V] {
	return KeysValues(next, convSlice(as.ErrTail(keyExtractor)), as.ErrTail(valsExtractor))
}

// KeyValuess transforms a loop to the key/value loop based on applying key, value extractor to the elements
func KeyValuess[T, K, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valsExtractor func(T) ([]V, error)) breakkvloop.Loop[K, V] {
	return KeysValues(next, convSlice(keyExtractor), valsExtractor)
}

// ExtraVals transforms a loop to the key/value loop based on applying value extractor to the elements
func ExtraVals[T, V any](next func() (T, bool, error), valsExtractor func(T) []V) breakkvloop.Loop[T, V] {
	return KeyValues(next, as.Is[T], valsExtractor)
}

// ExtraValss transforms a loop to the key/value loop based on applying values extractor to the elements
func ExtraValss[T, V any](next func() (T, bool, error), valsExtractor func(T) ([]V, error)) breakkvloop.Loop[T, V] {
	return KeyValuess(next, as.ErrTail(as.Is[T]), valsExtractor)
}

// ExtraKeys transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKeys[T, K any](next func() (T, bool, error), keysExtractor func(T) []K) breakkvloop.Loop[K, T] {
	return KeysValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyss transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKeyss[T, K any](next func() (T, bool, error), keyExtractor func(T) (K, error)) breakkvloop.Loop[K, T] {
	return KeyValuess(next, keyExtractor, as.ErrTail(convert.AsSlice[T]))
}

// ExtraKey transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKey[T, K any](next func() (T, bool, error), keysExtractor func(T) K) breakkvloop.Loop[K, T] {
	return KeyValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyy transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKeyy[T, K any](next func() (T, bool, error), keyExtractor func(T) (K, error)) breakkvloop.Loop[K, T] {
	return KeyValuee[T, K](next, keyExtractor, as.ErrTail(as.Is[T]))
}

// ExtraValue transforms a loop to the key/value loop based on applying value extractor to the elements
func ExtraValue[T, V any](next func() (T, bool, error), valueExtractor func(T) V) breakkvloop.Loop[T, V] {
	return KeyValue(next, as.Is[T], valueExtractor)
}

// ExtraValuee transforms a loop to the key/value loop based on applying value extractor to the elements
func ExtraValuee[T, V any](next func() (T, bool, error), valExtractor func(T) (V, error)) breakkvloop.Loop[T, V] {
	return KeyValuee[T, T, V](next, as.ErrTail(as.Is[T]), valExtractor)
}

// Group converts elements retrieved by the 'next' function into a new map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Group[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) K, valExtractor func(T) V) (map[K][]V, error) {
	return Groupp(next, as.ErrTail(keyExtractor), as.ErrTail(valExtractor))
}

// Groupp converts elements retrieved by the 'next' function into a new map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Groupp[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K][]V, error) {
	return MapResolvv(next, keyExtractor, valExtractor, func(ok bool, k K, rv []V, v V) ([]V, error) {
		return resolv.Slice(ok, k, rv, v), nil
	})
}

// GroupByMultiple converts elements retrieved by the 'next' function into a new map, extracting multiple keys, values per each element applying the 'keysExtractor' and 'valsExtractor' functions.
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

// GroupByMultipleKeys converts elements retrieved by the 'next' function into a new map, extracting multiple keys, one value per each element applying the 'keysExtractor' and 'valExtractor' functions.
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

// GroupByMultipleValues converts elements retrieved by the 'next' function into a new map, extracting one key, multiple values per each element applying the 'keyExtractor' and 'valsExtractor' functions.
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

// Map collects key\value elements into a new map by iterating over the elements
func Map[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) (map[K]V, error) {
	return Mapp(From(next), as.ErrTail(keyExtractor), as.ErrTail(valExtractor))
}

// Mapp collects key\value elements into a new map by iterating over the elements
func Mapp[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K]V, error) {
	return MapResolvv(next, keyExtractor, valExtractor, func(ok bool, k K, rv V, v V) (V, error) { return resolv.First(ok, k, rv, v), nil })
}

// MapResolvv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values
func MapResolvv[T any, K comparable, V, VR any](
	next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error),
	resolver func(bool, K, VR, V) (VR, error),
) (m map[K]VR, err error) {
	if next == nil {
		return nil, nil
	}
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
func ConvertAndReduce[From, To any](next func() (From, bool, error), converter func(From) To, merge func(To, To) To) (out To, err error) {
	return Reduce(Convert(next, converter), merge)
}

// ConvAndReduce converts each elements and merges them into one
func ConvAndReduce[From, To any](next func() (From, bool, error), converter func(From) (To, error), merge func(To, To) To) (out To, err error) {
	return Reduce(Conv(next, converter), merge)
}

// Crank rertieves a next element from the 'next' function, returns the function, element, successfully flag.
func Crank[T any](next func() (T, bool, error)) (n Loop[T], t T, ok bool, err error) {
	if next != nil {
		t, ok, err = next()
	}
	return next, t, ok, err
}

func brk(err error) error {
	if errors.Is(err, c.Break) {
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
