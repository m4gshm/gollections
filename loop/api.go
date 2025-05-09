// Package loop provides helpers for loop operation and iterator implementations
//
// Deprecated: use the [github.com/m4gshm/gollections/seq], [github.com/m4gshm/gollections/seqe], [github.com/m4gshm/gollections/seq2] packages API instead.
package loop

import (
	"unsafe"

	"golang.org/x/exp/constraints"

	breakkvloop "github.com/m4gshm/gollections/break/kv/loop"
	breakloop "github.com/m4gshm/gollections/break/loop"
	breakAlways "github.com/m4gshm/gollections/break/predicate/always"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/predicate/always"
)

// Break is the 'break' statement of the For, Track methods.
var Break = c.Break

// Continue is an alias of the nil value used to continue iterating by For, Track methods.
var Continue = c.Continue

// S wrap the elements by loop function.
func S[TS ~[]T, T any](elements TS) Loop[T] {
	return Of(elements...)
}

// Of wrap the elements by loop function.
func Of[T any](elements ...T) Loop[T] {
	l := len(elements)
	if l == 0 {
		return nil
	}
	i := 0
	return func() (T, bool) {
		if i < l {
			e, ok := elements[i], true
			i++
			return e, ok
		}
		var e T
		return e, false
	}
}

// All is an adapter for the next function for iterating by `for ... range`.
func All[T any](next func() (T, bool), consumer func(T) bool) {
	if next == nil {
		return
	}
	for v, ok := next(); ok && consumer(v); v, ok = next() {
	}
}

// New makes a loop from an abstract source
func New[S, T any](source S, hasNext func(S) bool, getNext func(S) T) Loop[T] {
	return func() (out T, ok bool) {
		if hasNext(source) {
			out, ok = getNext(source), true
		}
		return out, ok
	}
}

// For applies the 'consumer' function for the elements retrieved by the 'next' function until the consumer returns the c.Break to stop.
func For[T any](next func() (T, bool), consumer func(T) error) error {
	if next == nil {
		return nil
	}
	for v, ok := next(); ok; v, ok = next() {
		if err := consumer(v); err == Break {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies the 'consumer' function to the elements retrieved by the 'next' function
func ForEach[T any](next func() (T, bool), consumer func(T)) {
	if next == nil {
		return
	}
	for v, ok := next(); ok; v, ok = next() {
		consumer(v)
	}
}

// ForEachFiltered applies the 'consumer' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func ForEachFiltered[T any](next func() (T, bool), predicate func(T) bool, consumer func(T)) {
	if next == nil {
		return
	}
	for v, ok := next(); ok; v, ok = next() {
		if predicate(v) {
			consumer(v)
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool), predicate func(T) bool) (v T, ok bool) {
	if next == nil {
		return v, false
	}
	for one, ok := next(); ok; one, ok = next() {
		if predicate(one) {
			return one, true
		}
	}
	return v, ok
}

// Firstt returns the first element that satisfies the condition of the 'predicate' function
func Firstt[T any](next func() (T, bool), predicate func(T) (bool, error)) (v T, ok bool, err error) {
	if next == nil {
		return v, false, nil
	}
	for {
		if out, ok := next(); !ok {
			return out, false, nil
		} else if ok, err := predicate(out); err != nil || ok {
			return out, ok, err
		}
	}
}

// Track applies the 'consumer' function to position/element pairs retrieved by the 'next' function until the consumer returns the c.Break to stop.tracking.
func Track[I, T any](next func() (I, T, bool), consumer func(I, T) error) error {
	return kvloop.Track(next, consumer)
}

// TrackEach applies the 'consumer' function to position/element pairs retrieved by the 'next' function
func TrackEach[I, T any](next func() (I, T, bool), consumer func(I, T)) {
	kvloop.TrackEach(next, consumer)
}

// Slice collects the elements retrieved by the 'next' function into a new slice
func Slice[T any](next func() (T, bool)) []T {
	return SliceCap(next, 0)
}

// SliceCap collects the elements retrieved by the 'next' function into a new slice with predefined capacity
func SliceCap[T any](next func() (T, bool), capacity int) (out []T) {
	if next == nil {
		return nil
	}
	if capacity > 0 {
		out = make([]T, 0, capacity)
	}
	return Append(next, out)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func Append[T any, TS ~[]T](next func() (T, bool), out TS) TS {
	if next == nil {
		return out
	}
	for {
		v, ok := next()
		if !ok {
			break
		}
		out = append(out, v)
	}
	return out
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func Reduce[T any](next func() (T, bool), merge func(T, T) T) T {
	result, _ := ReduceOK(next, merge)
	return result
}

// ReduceOK reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceOK[T any](next func() (T, bool), merge func(T, T) T) (result T, ok bool) {
	if next == nil {
		return result, false
	}
	if result, ok = next(); !ok {
		return result, false
	}
	return Accum(result, next, merge), true
}

// Reducee reduces the elements retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func Reducee[T any](next func() (T, bool), merge func(T, T) (T, error)) (T, error) {
	result, _, err := ReduceeOK(next, merge)
	return result, err
}

// ReduceeOK reduces the elements retrieved by the 'next' function into an one pair using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func ReduceeOK[T any](next func() (T, bool), merge func(T, T) (T, error)) (result T, ok bool, err error) {
	if next == nil {
		return result, false, nil
	}
	if result, ok = next(); !ok {
		return result, false, nil
	}
	result, err = Accumm(result, next, merge)
	return result, true, err
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func Accum[T any](first T, next func() (T, bool), merge func(T, T) T) T {
	accumulator := first
	if next == nil {
		return accumulator
	}
	for v, ok := next(); ok; v, ok = next() {
		accumulator = merge(accumulator, v)
	}
	return accumulator
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func Accumm[T any](first T, next func() (T, bool), merge func(T, T) (T, error)) (accumulator T, err error) {
	accumulator = first
	if next == nil {
		return accumulator, nil
	}
	for v, ok := next(); ok; v, ok = next() {
		accumulator, err = merge(accumulator, v)
		if err != nil {
			return accumulator, err
		}
	}
	return accumulator, nil
}

// Sum returns the sum of all elements
func Sum[T c.Summable](next func() (T, bool)) (out T) {
	return Accum(out, next, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func HasAny[T any](next func() (T, bool), predicate func(T) bool) bool {
	_, ok := First(next, predicate)
	return ok
}

// Contains finds the first element that equal to the example and returns true
func Contains[T comparable](next func() (T, bool), example T) bool {
	if next == nil {
		return false
	}
	for one, ok := next(); ok; one, ok = next() {
		if one == example {
			return true
		}
	}
	return false
}

// Conv creates a loop that applies the 'converter' function to iterable elements.
func Conv[From, To any](next func() (From, bool), converter func(From) (To, error)) breakloop.Loop[To] {
	return breakloop.Conv(breakloop.From(next), converter)
}

// ConvS creates a loop that applies the 'converter' function to the 'elements' slice.
func ConvS[FS ~[]From, From, To any](elements FS, converter func(From) (To, error)) breakloop.Loop[To] {
	return Conv(S(elements), converter)
}

// Convert creates a loop that applies the 'converter' function to iterable elements.
func Convert[From, To any](next func() (From, bool), converter func(From) To) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool) {
		v, ok := next()
		if ok {
			return converter(v), true
		}
		return t, false
	}
}

// ConvertS creates a loop that applies the 'converter' function to the 'elements' slice.
func ConvertS[FS ~[]From, From, To any](elements FS, converter func(From) To) Loop[To] {
	return Convert(S(elements), converter)
}

// ConvOK creates a loop that applies the 'converter' function to iterable elements.
// The converter may returns a value or ok=false to exclude the value from the loop.
// It may also return an error to abort the loop.
func ConvOK[From, To any](next func() (From, bool), converter func(from From) (To, bool, error)) breakloop.Loop[To] {
	return breakloop.ConvOK(breakloop.From(next), converter)
}

// ConvertOK creates a loop that applies the 'converter' function to iterable elements.
// The converter may returns a value or ok=false to exclude the value from the loop.
func ConvertOK[From, To any](next func() (From, bool), converter func(from From) (To, bool)) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool) {
		for e, ok := next(); ok; e, ok = next() {
			if t, ok := converter(e); ok {
				return t, true
			}
		}
		return t, false
	}
}

// FiltAndConv creates a loop that filters source elements and converts them
func FiltAndConv[From, To any](next func() (From, bool), filter func(From) (bool, error), converter func(From) (To, error)) breakloop.Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool, err error) {
		for {
			if f, ok, err := Firstt(next, filter); err != nil || !ok {
				return t, false, err
			} else if cf, err := converter(f); err != nil {
				return t, false, err
			} else {
				return cf, true, nil
			}
		}
	}
}

// FilterAndConvert creates a loop that filters source elements and converts them
func FilterAndConvert[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To) Loop[To] {
	return FilterConvertFilter(next, filter, converter, always.True[To])
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To, filterTo func(To) bool) Loop[To] {
	if next == nil {
		return nil
	}
	return func() (t To, ok bool) {
		for {
			if f, ok := First(next, filter); !ok {
				return t, false
			} else if t = converter(f); filterTo(t) {
				return t, ok
			}
		}
	}
}

// ConvertAndFilter additionally filters 'To' elements
func ConvertAndFilter[From, To any](next func() (From, bool), converter func(From) To, filter func(To) bool) Loop[To] {
	return FilterConvertFilter(next, always.True[From], converter, filter)
}

// Flatt converts a two-dimensional loop in an one-dimensional one.
func Flatt[From, To any](next func() (From, bool), flattener func(From) ([]To, error)) breakloop.Loop[To] {
	return breakloop.Flatt(breakloop.From(next), flattener)
}

// FlattS transforms the n-dimensional 'elements' slice to a n-1 dimensional loop.
func FlattS[FS ~[]From, From, To any](elements FS, flattener func(From) ([]To, error)) breakloop.Loop[To] {
	return Flatt(S(elements), flattener)
}

// Flat converts a two-dimensional loop in a one-dimensional one, like:
//
//	var arrays func() ([]int, boot) = ...
//	var ints func() (int, boot) = loop.Flat(arrays, as.Is)
func Flat[From, To any](next func() (From, bool), flattener func(From) []To) Loop[To] {
	if next == nil {
		return nil
	}
	var (
		elemSizeTo      uintptr = notsafe.GetTypeSize[To]()
		arrayTo         unsafe.Pointer
		indexTo, sizeTo int
	)
	return func() (t To, ok bool) {
		if sizeTo > 0 {
			if indexTo < sizeTo {
				i := indexTo
				indexTo++
				return *(*To)(notsafe.GetArrayElemRef(arrayTo, i, elemSizeTo)), true
			}
			indexTo = 0
			arrayTo = nil
			sizeTo = 0
		}
		for {
			if v, ok := next(); !ok {
				var no To
				return no, false
			} else if elementsTo := flattener(v); len(elementsTo) > 0 {
				indexTo = 1
				header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
				arrayTo = unsafe.Pointer(header.Data)
				sizeTo = header.Len
				return *(*To)(notsafe.GetArrayElemRef(arrayTo, 0, elemSizeTo)), true
			}
		}
	}
}

// FlatS creates a loop that extracts slices of 'To' by the 'flattener' function from the elements of 'From' and flattens as one iterable collection of 'To' elements.
func FlatS[FS ~[]From, From, To any](elements FS, flattener func(From) []To) Loop[To] {
	return Flat(S(elements), flattener)
}

// FiltAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FiltAndFlat[From, To any](next func() (From, bool), filter func(From) (bool, error), flattener func(From) ([]To, error)) breakloop.Loop[To] {
	return breakloop.FiltFlattFilt(breakloop.From(next), filter, flattener, breakAlways.True[To])
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any](next func() (From, bool), filter func(From) bool, flattener func(From) []To) Loop[To] {
	return FilterFlatFilter(next, filter, flattener, always.True[To])
}

// FlatAndFilt extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlatAndFilt[From, To any](next func() (From, bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) breakloop.Loop[To] {
	return breakloop.FiltFlattFilt(next, breakAlways.True[From], flattener, filterTo)
}

// FlattAndFilter extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlattAndFilter[From, To any](next func() (From, bool), flattener func(From) []To, filterTo func(To) bool) Loop[To] {
	return FilterFlatFilter(next, always.True[From], flattener, filterTo)
}

// FiltFlattFilt filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FiltFlattFilt[From, To any](next func() (From, bool), filterFrom func(From) (bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) breakloop.Loop[To] {
	return breakloop.FiltFlattFilt(breakloop.From(next), filterFrom, flattener, filterTo)
}

// FilterFlatFilter filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FilterFlatFilter[From, To any](next func() (From, bool), filterFrom func(From) bool, flattener func(From) []To, filterTo func(To) bool) Loop[To] {
	if next == nil {
		return nil
	}
	var (
		elemSizeTo      uintptr = notsafe.GetTypeSize[To]()
		arrayTo         unsafe.Pointer
		indexTo, sizeTo int
	)
	return func() (t To, ok bool) {
		for {
			if sizeTo > 0 {
				if indexTo < sizeTo {
					i := indexTo
					indexTo++
					t = *(*To)(notsafe.GetArrayElemRef(arrayTo, i, elemSizeTo))
					if ok := filterTo(t); ok {
						return t, true
					}
				}
				indexTo = 0
				arrayTo = nil
				sizeTo = 0
			}

			if v, ok := next(); !ok {
				return t, false
			} else if filterFrom(v) {
				if elementsTo := flattener(v); len(elementsTo) > 0 {
					indexTo = 1
					header := notsafe.GetSliceHeaderByRef(unsafe.Pointer(&elementsTo))
					arrayTo = unsafe.Pointer(header.Data)
					sizeTo = header.Len
					t = *(*To)(notsafe.GetArrayElemRef(arrayTo, 0, elemSizeTo))
					if ok := filterTo(t); ok {
						return t, true
					}
				}
			}
		}
	}
}

// Filt creates a loop that checks elements by the 'filter' function and returns successful ones.
func Filt[T any](next func() (T, bool), filter func(T) (bool, error)) breakloop.Loop[T] {
	return breakloop.Filt(breakloop.From(next), filter)
}

// FiltS creates a loop that checks slice elements by the 'filter' function and returns successful ones.
func FiltS[TS ~[]T, T any](elements TS, filter func(T) (bool, error)) breakloop.Loop[T] {
	return Filt(S(elements), filter)
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func Filter[T any](next func() (T, bool), filter func(T) bool) Loop[T] {
	if next == nil {
		return nil
	}
	return func() (T, bool) {
		return First(next, filter)
	}
}

// FilterS creates a loop that checks slice elements by the 'filter' function and returns successful ones.
func FilterS[TS ~[]T, T any](elements TS, filter func(T) bool) Loop[T] {
	return Filter(S(elements), filter)
}

// NotNil creates a loop that filters nullable elements
func NotNil[T any](next func() (*T, bool)) Loop[*T] {
	return Filter(next, not.Nil[T])
}

// PtrVal creates a loop that transform pointers to the values referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any](next func() (*T, bool)) Loop[T] {
	return Convert(next, convert.PtrVal[T])
}

// NoNilPtrVal creates a loop that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any](next func() (*T, bool)) Loop[T] {
	return ConvertOK(next, convert.NoNilPtrVal[T])
}

// KeyValue transforms a loop to the key/value loop based on applying key, value extractors to the elements
func KeyValue[T any, K, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) kvloop.Loop[K, V] {
	if next == nil {
		return nil
	}
	return func() (key K, value V, ok bool) {
		if elem, nextOk := next(); nextOk {
			key = keyExtractor(elem)
			value = valExtractor(elem)
			ok = true
		}
		return key, value, ok
	}
}

// KeyValuee transforms a loop to the key/value loop based on applying key, value extractors to the elements
func KeyValuee[T any, K, V any](next func() (T, bool), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) breakkvloop.Loop[K, V] {
	return breakloop.KeyValuee(breakloop.From(next), keyExtractor, valExtractor)
}

// KeysValues transforms a loop to the key/value loop based on applying multiple keys, values extractor to the elements
func KeysValues[T, K, V any](next func() (T, bool), keysExtractor func(T) []K, valsExtractor func(T) []V) kvloop.Loop[K, V] {
	if next == nil {
		return nil
	}
	var (
		keys   []K
		values []V
		ki, vi int
	)
	return func() (key K, value V, ok bool) {
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
			} else if elem, nextOk := next(); nextOk {
				keys = keysExtractor(elem)
				values = valsExtractor(elem)
				ki, vi = 0, 0
			} else {
				keys, values = nil, nil
				break
			}
		}
		return key, value, ok
	}
}

// KeysValue transforms a loop to the key/value loop based on applying keys, value extractor to the elements
func KeysValue[T, K, V any](next func() (T, bool), keysExtractor func(T) []K, valExtractor func(T) V) kvloop.Loop[K, V] {
	return KeysValues(next, keysExtractor, func(t T) []V { return convert.AsSlice(valExtractor(t)) })
}

// KeysValuee transforms a loop to the key/value loop based on applying keys, value extractor to the elements
func KeysValuee[T, K, V any](next func() (T, bool), keysExtractor func(T) ([]K, error), valExtractor func(T) (V, error)) breakkvloop.Loop[K, V] {
	return breakloop.KeysValuee(breakloop.From(next), keysExtractor, valExtractor)
}

// KeyValues transforms a loop to the key/value loop based on applying key, values extractor to the elements
func KeyValues[T, K, V any](next func() (T, bool), keyExtractor func(T) K, valsExtractor func(T) []V) kvloop.Loop[K, V] {
	return KeysValues(next, func(t T) []K { return convert.AsSlice(keyExtractor(t)) }, valsExtractor)
}

// KeyValuess transforms a loop to the key/value loop based on applying key, values extractor to the elements
func KeyValuess[T, K, V any](next func() (T, bool), keyExtractor func(T) (K, error), valsExtractor func(T) ([]V, error)) breakkvloop.Loop[K, V] {
	return breakloop.KeyValuess(breakloop.From(next), keyExtractor, valsExtractor)
}

// ExtraVals transforms a loop to the key/value loop based on applying values extractor to the elements
func ExtraVals[T, V any](next func() (T, bool), valsExtractor func(T) []V) kvloop.Loop[T, V] {
	return KeyValues(next, as.Is[T], valsExtractor)
}

// ExtraValss transforms a loop to the key/value loop based on applying values extractor to the elements
func ExtraValss[T, V any](next func() (T, bool), valsExtractor func(T) ([]V, error)) breakkvloop.Loop[T, V] {
	return KeyValuess(next, as.ErrTail(as.Is[T]), valsExtractor)
}

// ExtraKeys transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKeys[T, K any](next func() (T, bool), keysExtractor func(T) []K) kvloop.Loop[K, T] {
	return KeysValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyss transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKeyss[T, K any](next func() (T, bool), keyExtractor func(T) (K, error)) breakkvloop.Loop[K, T] {
	return KeyValuess(next, keyExtractor, as.ErrTail(convert.AsSlice[T]))
}

// ExtraKey transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKey[T, K any](next func() (T, bool), keysExtractor func(T) K) kvloop.Loop[K, T] {
	return KeyValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyy transforms a loop to the key/value loop based on applying key extractor to the elements
func ExtraKeyy[T, K any](next func() (T, bool), keyExtractor func(T) (K, error)) breakkvloop.Loop[K, T] {
	return breakloop.KeyValuee[T, K](breakloop.From(next), keyExtractor, as.ErrTail(as.Is[T]))
}

// ExtraValue transforms a loop to the key/value loop based on applying value extractor to the elements
func ExtraValue[T, V any](next func() (T, bool), valueExtractor func(T) V) kvloop.Loop[T, V] {
	return KeyValue(next, as.Is[T], valueExtractor)
}

// ExtraValuee transforms a loop to the key/value loop based on applying value extractor to the elements
func ExtraValuee[T, V any](next func() (T, bool), valExtractor func(T) (V, error)) breakkvloop.Loop[T, V] {
	return breakloop.KeyValuee[T, T, V](breakloop.From(next), as.ErrTail(as.Is[T]), valExtractor)
}

// Group converts elements retrieved by the 'next' function into a new map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Group[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	return MapResolv(next, keyExtractor, valExtractor, resolv.Slice[K, V])
}

// Groupp converts elements retrieved by the 'next' function into a new map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to a value.
func Groupp[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K][]V, error) {
	return breakloop.Groupp(breakloop.From(next), keyExtractor, valExtractor)
}

// GroupByMultiple converts elements retrieved by the 'next' function into a new map, extracting multiple keys, values per each element applying the 'keysExtractor' and 'valsExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valsExtractor retrieves one or more values per element.
func GroupByMultiple[T any, K comparable, V any](next func() (T, bool), keysExtractor func(T) []K, valsExtractor func(T) []V) map[K][]V {
	if next == nil {
		return nil
	}
	groups := map[K][]V{}
	for e, ok := next(); ok; e, ok = next() {
		if keys, vals := keysExtractor(e), valsExtractor(e); len(keys) == 0 {
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
	return groups
}

// GroupByMultipleKeys converts elements retrieved by the 'next' function into a new map, extracting multiple keys, one value per each element applying the 'keysExtractor' and 'valExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valExtractor converts an element to a value.
func GroupByMultipleKeys[T any, K comparable, V any](next func() (T, bool), keysExtractor func(T) []K, valExtractor func(T) V) map[K][]V {
	if next == nil {
		return nil
	}
	groups := map[K][]V{}
	for e, ok := next(); ok; e, ok = next() {
		if keys, v := keysExtractor(e), valExtractor(e); len(keys) == 0 {
			var key K
			initGroup(key, v, groups)
		} else {
			for _, key := range keys {
				initGroup(key, v, groups)
			}
		}
	}
	return groups
}

// GroupByMultipleValues converts elements retrieved by the 'next' function into a new map, extracting one key, multiple values per each element applying the 'keyExtractor' and 'valsExtractor' functions.
// The keyExtractor converts an element to a key.
// The valsExtractor retrieves one or more values per element.
func GroupByMultipleValues[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valsExtractor func(T) []V) map[K][]V {
	if next == nil {
		return nil
	}
	groups := map[K][]V{}
	for e, ok := next(); ok; e, ok = next() {
		if key, vals := keyExtractor(e), valsExtractor(e); len(vals) == 0 {
			var v V
			initGroup(key, v, groups)
		} else {
			for _, v := range vals {
				initGroup(key, v, groups)
			}
		}
	}
	return groups
}

func initGroup[T any, K comparable, TS ~[]T](key K, e T, groups map[K]TS) {
	groups[key] = append(groups[key], e)
}

// Map collects key\value elements into a new map by iterating over the elements
func Map[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) map[K]V {
	return MapResolv(next, keyExtractor, valExtractor, resolv.First[K, V])
}

// Mapp collects key\value elements into a new map by iterating over the elements
func Mapp[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K]V, error) {
	return breakloop.Mapp(breakloop.From(next), keyExtractor, valExtractor)
}

// MapResolv collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values
func MapResolv[T any, K comparable, V, VR any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR) map[K]VR {
	return AppendMapResolv(next, keyExtractor, valExtractor, resolver, nil)
}

// MapResolvOrder collects key\value elements into a new map by iterating over the elements with resolving of duplicated key values.
// Returns a slice with the keys ordered by the time they were added and the resolved key\value map.
func MapResolvOrder[TS ~[]T, T any, K comparable, V, VR any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR) ([]K, map[K]VR) {
	return AppendMapResolvOrder(next, keyExtractor, valExtractor, resolver, nil, nil)
}

// AppendMapResolv collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values
func AppendMapResolv[T any, K comparable, V, VR any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR, dest map[K]VR) map[K]VR {
	if next == nil {
		return nil
	}
	if dest == nil {
		dest = map[K]VR{}
	}
	for e, ok := next(); ok; e, ok = next() {
		k, v := keyExtractor(e), valExtractor(e)
		exists, ok := dest[k]
		dest[k] = resolver(ok, k, exists, v)
	}
	return dest
}

// AppendMapResolvOrder collects key\value elements into the 'dest' map by iterating over the elements with resolving of duplicated key values
// Additionaly populates the 'order' slice by the keys ordered by the time they were added and the resolved key\value map.
func AppendMapResolvOrder[T any, K comparable, V, VR any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR, order []K, dest map[K]VR) ([]K, map[K]VR) {
	if next == nil {
		return nil, nil
	}
	if dest == nil {
		dest = map[K]VR{}
	}
	for e, ok := next(); ok; e, ok = next() {
		k, v := keyExtractor(e), valExtractor(e)
		exists, ok := dest[k]
		dest[k] = resolver(ok, k, exists, v)
		if !ok {
			order = append(order, k)
		}
	}
	return order, dest
}

// Series makes a sequence by applying the 'next' function to the previous step generated value.
func Series[T any](first T, next func(T) (T, bool)) Loop[T] {
	if next == nil {
		return nil
	}
	current := first
	init := true
	return func() (out T, ok bool) {
		if init {
			init = false
			return current, true
		} else {
			next, ok := next(current)
			current = next
			return current, ok
		}
	}
}

// RangeClosed creates a loop that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer | rune](from T, toInclusive T) Loop[T] {
	amount := toInclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	amount++
	nextElement := from
	i := T(0)
	return func() (out T, ok bool) {
		if ok = i < amount; ok {
			out = nextElement
			i++
			nextElement = nextElement + delta
		}
		return out, ok
	}
}

// Range creates a loop that generates integers in the range defined by from and to exclusive
func Range[T constraints.Integer | rune](from T, toExclusive T) Loop[T] {
	amount := toExclusive - from
	delta := T(1)
	if amount < 0 {
		amount = -amount
		delta = -delta
	}
	nextElement := from
	i := T(0)
	return func() (out T, ok bool) {
		if ok = i < amount; ok {
			out = nextElement
			i++
			nextElement = nextElement + delta
		}
		return out, ok
	}
}

// OfIndexed builds a loop by extracting elements from an indexed soruce.
// the len is length ot the source.
// the getAt retrieves an element by its index from the source.
func OfIndexed[T any](amount int, getAt func(int) T) Loop[T] {
	if getAt == nil {
		return nil
	}
	i := 0
	return func() (out T, ok bool) {
		if ok = i < amount; ok {
			out = getAt(i)
			i++
		}
		return out, ok
	}
}

// ConvertAndReduce converts each elements and merges them into one
func ConvertAndReduce[From, To any](next func() (From, bool), converter func(From) To, merge func(To, To) To) (out To) {
	if next == nil {
		return out
	}
	if v, ok := next(); ok {
		out = converter(v)
	} else {
		return out
	}
	for v, ok := next(); ok; v, ok = next() {
		out = merge(out, converter(v))
	}
	return out
}

// ConvAndReduce converts each elements and merges them into one
func ConvAndReduce[From, To any](next func() (From, bool), converter func(From) (To, error), merge func(To, To) To) (out To, err error) {
	if next == nil {
		return out, nil
	}
	if v, ok := next(); ok {
		out, err = converter(v)
		if err != nil {
			return out, err
		}
	} else {
		return out, nil
	}
	for v, ok := next(); ok; v, ok = next() {
		c, err := converter(v)
		if err != nil {
			return out, err
		}
		out = merge(out, c)
	}
	return out, nil
}

// Crank rertieves a next element from the 'next' function, returns the function, element, successfully flag.
func Crank[T any](next func() (T, bool)) (n Loop[T], t T, ok bool) {
	if next != nil {
		t, ok = next()
	}
	return next, t, ok
}
