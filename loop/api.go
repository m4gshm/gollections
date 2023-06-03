// Package loop provides helpers for loop operation and iterator implementations
package loop

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/break/loop"
	breakAlways "github.com/m4gshm/gollections/break/predicate/always"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/predicate/always"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = c.ErrBreak

// Looper provides an iterable loop function
type Looper[T any, I interface{ Next() (T, bool) }] interface {
	Loop() I
}

// Of wrap the elements by loop function
func Of[T any](elements ...T) func() (e T, ok bool) {
	l := len(elements)
	i := 0
	if l == 0 || i < 0 || i >= l {
		return func() (e T, ok bool) { return e, false }
	}
	return func() (e T, ok bool) {
		if i < l {
			e, ok = elements[i], true
			i++
		}
		return e, ok
	}
}

type F struct {
	stop bool
}

func Break(f *F) {
	f.stop = true
}

func ForB[T any](next func() (T, bool), walker func(*F, T)) {
	f := &F{}
	for v, ok := next(); ok; v, ok = next() {
		walker(f, v)
		if f.stop {
			break
		}
	}
}

// For applies the 'walker' function for the elements retrieved by the 'next' function. Return the c.ErrBreak to stop
func For[T any](next func() (T, bool), walker func(T) error) error {
	for v, ok := next(); ok; v, ok = next() {
		if err := walker(v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies the 'walker' function to the elements retrieved by the 'next' function
func ForEach[T any](next func() (T, bool), walker func(T)) {
	for v, ok := next(); ok; v, ok = next() {
		walker(v)
	}
}

// ForEachFiltered applies the 'walker' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func ForEachFiltered[T any](next func() (T, bool), predicate func(T) bool, walker func(T)) {
	for v, ok := next(); ok; v, ok = next() {
		if predicate(v) {
			walker(v)
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool), predicate func(T) bool) (v T, ok bool) {
	for one, ok := next(); ok; one, ok = next() {
		if predicate(one) {
			return one, true
		}
	}
	return v, ok
}

// Track applies the 'tracker' function to position/element pairs retrieved by the 'next' function. Return the c.ErrBreak to stop tracking..
func Track[I, T any](next func() (I, T, bool), tracker func(I, T) error) error {
	for p, v, ok := next(); ok; p, v, ok = next() {
		if err := tracker(p, v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEach applies the 'tracker' function to position/element pairs retrieved by the 'next' function
func TrackEach[I, T any](next func() (I, T, bool), tracker func(I, T)) {
	for p, v, ok := next(); ok; p, v, ok = next() {
		tracker(p, v)
	}
}

// Slice collects the elements retrieved by the 'next' function into a new slice
func Slice[T any](next func() (T, bool)) []T {
	return SliceCap(next, 0)
}

// SliceCap collects the elements retrieved by the 'next' function into a new slice with predefined capacity
func SliceCap[T any](next func() (T, bool), cap int) (out []T) {
	if cap > 0 {
		out = make([]T, 0, cap)
	}
	return Append(next, out)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func Append[T any, TS ~[]T](next func() (T, bool), out TS) TS {
	for v, ok := next(); ok; v, ok = next() {
		out = append(out, v)
	}
	return out
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merger' function
func Reduce[T any](next func() (T, bool), merger func(T, T) T) (result T) {
	if v, ok := next(); ok {
		result = v
	} else {
		return result
	}
	for v, ok := next(); ok; v, ok = next() {
		result = merger(result, v)
	}
	return result
}

// Sum returns the sum of all elements
func Sum[T c.Summable](next func() (T, bool)) T {
	return Reduce(next, op.Sum[T])
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func HasAny[T any](next func() (T, bool), predicate func(T) bool) bool {
	_, ok := First(next, predicate)
	return ok
}

// Contains  finds the first element that equal to the example and returns true
func Contains[T comparable](next func() (T, bool), example T) bool {
	for one, ok := next(); ok; one, ok = next() {
		if one == example {
			return true
		}
	}
	return false
}

// Conv instantiates an iterator that converts elements with a converter and returns them.
func Conv[From, To any](next func() (From, bool), converter func(From) (To, error)) loop.ConvertIter[From, To] {
	return loop.Conv(loop.From(next), converter)
}

// Convert instantiates an iterator that converts elements with a converter and returns them.
func Convert[From, To any](next func() (From, bool), converter func(From) To) ConvertIter[From, To] {
	return ConvertIter[From, To]{next: next, converter: converter}
}

// ConvCheck is similar to ConvertFilt, but it checks and transforms elements together
func ConvCheck[From, To any](next func() (From, bool), converter func(from From) (To, bool, error)) loop.ConvertCheckIter[From, To] {
	return loop.ConvCheck(loop.From(next), converter)
}

// ConvertCheck is similar to ConvertFilt, but it checks and transforms elements together
func ConvertCheck[From, To any](next func() (From, bool), converter func(from From) (To, bool)) ConvertCheckIter[From, To] {
	return ConvertCheckIter[From, To]{next: next, converter: converter}
}

// FiltAndConv returns a stream that filters source elements and converts them
func FiltAndConv[From, To any](next func() (From, bool), filter func(From) (bool, error), converter func(From) (To, error)) loop.ConvFiltIter[From, To] {
	return loop.FiltAndConv(loop.From(next), filter, converter)
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To) ConvertFiltIter[From, To] {
	return FilterConvertFilter(next, filter, converter, always.True[To])
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To, filterTo func(To) bool) ConvertFiltIter[From, To] {
	return ConvertFiltIter[From, To]{next: next, converter: converter, filterFrom: filter, filterTo: filterTo}
}

// ConvertAndFilter additionally filters 'To' elements
func ConvertAndFilter[From, To any](next func() (From, bool), converter func(From) To, filter func(To) bool) ConvertFiltIter[From, To] {
	return FilterConvertFilter(next, always.True[From], converter, filter)
}

// Flatt instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](next func() (From, bool), flattener func(From) ([]To, error)) *loop.FlatIter[From, To] {
	return loop.Flatt(loop.From(next), flattener)
}

// Flat instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flat[From, To any](next func() (From, bool), flattener func(From) []To) *FlatIter[From, To] {
	return &FlatIter[From, To]{next: next, flattener: flattener, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FiltAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FiltAndFlat[From, To any](next func() (From, bool), filter func(From) (bool, error), flattener func(From) ([]To, error)) *loop.FlattFiltIter[From, To] {
	return loop.FiltFlattFilt(loop.From(next), filter, flattener, breakAlways.True[To])
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[From, To any](next func() (From, bool), filter func(From) bool, flattener func(From) []To) *FlatFilterIter[From, To] {
	return FilterFlatFilter(next, filter, flattener, always.True[To])
}

// FlatAndFilt extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlatAndFilt[From, To any](next func() (From, bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) *loop.FlattFiltIter[From, To] {
	return loop.FiltFlattFilt(next, breakAlways.True[From], flattener, filterTo)
}

// FlattAndFilter extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlattAndFilter[From, To any](next func() (From, bool), flattener func(From) []To, filterTo func(To) bool) *FlatFilterIter[From, To] {
	return FilterFlatFilter(next, always.True[From], flattener, filterTo)
}

// FiltFlattFilt filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FiltFlattFilt[From, To any](next func() (From, bool), filterFrom func(From) (bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) *loop.FlattFiltIter[From, To] {
	return loop.FiltFlattFilt(loop.From(next), filterFrom, flattener, filterTo)
}

// FilterFlatFilter filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FilterFlatFilter[From, To any](next func() (From, bool), filterFrom func(From) bool, flattener func(From) []To, filterTo func(To) bool) *FlatFilterIter[From, To] {
	return &FlatFilterIter[From, To]{next: next, filterFrom: filterFrom, flattener: flattener, filterTo: filterTo, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Filt creates an iterator that checks elements by the 'filter' function and returns successful ones.
func Filt[T any](next func() (T, bool), filter func(T) (bool, error)) loop.FiltIter[T] {
	return loop.Filt(loop.From(next), filter)
}

// Filter creates an iterator that checks elements by the 'filter' function and returns successful ones.
func Filter[T any](next func() (T, bool), filter func(T) bool) FiltIter[T] {
	return FiltIter[T]{next: next, by: filter}
}

// NotNil creates an iterator that filters nullable elements
func NotNil[T any](next func() (*T, bool)) FiltIter[*T] {
	return Filter(next, not.Nil[T])
}

// PtrVal creates an iterator that transform pointers to the values referenced by those pointers.
// Nil pointers are transformet to zero values.
func PtrVal[T any](next func() (*T, bool)) ConvertIter[*T, T] {
	return Convert(next, convert.PtrVal[T])
}

// NoNilPtrVal creates an iterator that transform only not nil pointers to the values referenced referenced by those pointers.
// Nil pointers are ignored.
func NoNilPtrVal[T any](next func() (*T, bool)) ConvertCheckIter[*T, T] {
	return ConvertCheck(next, convert.NoNilPtrVal[T])
}

// KeyValue transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValue[T any, K, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) KeyValuer[T, K, V] {
	return NewKeyValuer(next, keyExtractor, valExtractor)
}

// KeyValuee transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func KeyValuee[T any, K, V any](next func() (T, bool), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) loop.KeyValuer[T, K, V] {
	return loop.KeyValuee(loop.From(next), keyExtractor, valExtractor)
}

// KeysValues transforms iterable elements to key/value iterator based on applying multiple keys, values extractor to the elements
func KeysValues[T, K, V any](next func() (T, bool), keysExtractor func(T) []K, valsExtractor func(T) []V) *MultipleKeyValuer[T, K, V] {
	return NewMultipleKeyValuer(next, keysExtractor, valsExtractor)
}

// KeysValue transforms iterable elements to key/value iterator based on applying keys, value extractor to the elements
func KeysValue[T, K, V any](next func() (T, bool), keysExtractor func(T) []K, valExtractor func(T) V) *MultipleKeyValuer[T, K, V] {
	return KeysValues(next, keysExtractor, func(t T) []V { return convert.AsSlice(valExtractor(t)) })
}

// KeysValuee transforms iterable elements to key/value iterator based on applying keys, value extractor to the elements
func KeysValuee[T, K, V any](next func() (T, bool), keysExtractor func(T) ([]K, error), valExtractor func(T) (V, error)) *loop.MultipleKeyValuer[T, K, V] {
	return loop.KeysValuee(loop.From(next), keysExtractor, valExtractor)
}

// KeyValues transforms iterable elements to key/value iterator based on applying key, values extractor to the elements
func KeyValues[T, K, V any](next func() (T, bool), keyExtractor func(T) K, valsExtractor func(T) []V) *MultipleKeyValuer[T, K, V] {
	return KeysValues(next, func(t T) []K { return convert.AsSlice(keyExtractor(t)) }, valsExtractor)
}

// KeyValuess transforms iterable elements to key/value iterator based on applying key, values extractor to the elements
func KeyValuess[T, K, V any](next func() (T, bool), keyExtractor func(T) (K, error), valsExtractor func(T) ([]V, error)) *loop.MultipleKeyValuer[T, K, V] {
	return loop.KeyValuess(loop.From(next), keyExtractor, valsExtractor)
}

// ExtraValues transforms iterable elements to key/value iterator based on applying values extractor to the elements
func ExtraValues[T, V any](next func() (T, bool), valsExtractor func(T) []V) *MultipleKeyValuer[T, T, V] {
	return KeyValues(next, as.Is[T], valsExtractor)
}

// ExtractValues transforms iterable elements to key/value iterator based on applying values extractor to the elements
func ExtraValuess[T, V any](next func() (T, bool), valsExtractor func(T) ([]V, error)) *loop.MultipleKeyValuer[T, T, V] {
	return KeyValuess(next, as.ErrTail(as.Is[T]), valsExtractor)
}

// ExtraKeys transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeys[T, K any](next func() (T, bool), keysExtractor func(T) []K) *MultipleKeyValuer[T, K, T] {
	return KeysValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyss transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeyss[T, K any](next func() (T, bool), keyExtractor func(T) (K, error)) *loop.MultipleKeyValuer[T, K, T] {
	return KeyValuess(next, keyExtractor, as.ErrTail(convert.AsSlice[T]))
}

// ExtraKey transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKey[T, K any](next func() (T, bool), keysExtractor func(T) K) KeyValuer[T, K, T] {
	return KeyValue(next, keysExtractor, as.Is[T])
}

// ExtraKeyy transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ExtraKeyy[T, K any](next func() (T, bool), keyExtractor func(T) (K, error)) loop.KeyValuer[T, K, T] {
	return loop.KeyValuee[T, K](loop.From(next), keyExtractor, as.ErrTail(as.Is[T]))
}

// ExtraValue transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraValue[T, V any](next func() (T, bool), valueExtractor func(T) V) KeyValuer[T, T, V] {
	return KeyValue(next, as.Is[T], valueExtractor)
}

// ExtraValue transforms iterable elements to key/value iterator based on applying value extractor to the elements
func ExtraValuee[T, V any](next func() (T, bool), valExtractor func(T) (V, error)) loop.KeyValuer[T, T, V] {
	return loop.KeyValuee[T, T, V](loop.From(next), as.ErrTail(as.Is[T]), valExtractor)
}

// Group converts elements retrieved by the 'next' function into a map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to an value.
func Group[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	return ToMapResolv(next, keyExtractor, valExtractor, resolv.Append[K, V])
}

// Groupp converts elements retrieved by the 'next' function into a map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to an value.
func Groupp[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K][]V, error) {
	return loop.Groupp(loop.From(next), keyExtractor, valExtractor)
}

// GroupByMultiple converts elements retrieved by the 'next' function into a map, extracting multiple keys, values per each element applying the 'keysExtractor' and 'valsExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valsExtractor retrieves one or more values per element.
func GroupByMultiple[T any, K comparable, V any](next func() (T, bool), keysExtractor func(T) []K, valsExtractor func(T) []V) map[K][]V {
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

// GroupByMultipleKeys converts elements retrieved by the 'next' function into a map, extracting multiple keys, one value per each element applying the 'keysExtractor' and 'valExtractor' functions.
// The keysExtractor retrieves one or more keys per element.
// The valExtractor converts an element to a value.
func GroupByMultipleKeys[T any, K comparable, V any](next func() (T, bool), keysExtractor func(T) []K, valExtractor func(T) V) map[K][]V {
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

// GroupByMultipleValues converts elements retrieved by the 'next' function into a map, extracting one key, multiple values per each element applying the 'keyExtractor' and 'valsExtractor' functions.
// The keyExtractor converts an element to a key.
// The valsExtractor retrieves one or more values per element.
func GroupByMultipleValues[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valsExtractor func(T) []V) map[K][]V {
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

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) map[K]V {
	return ToMapResolv(next, keyExtractor, valExtractor, resolv.First[K, V])
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[T any, K comparable, V, VR any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR) map[K]VR {
	m := map[K]VR{}
	for e, ok := next(); ok; e, ok = next() {
		k, v := keyExtractor(e), valExtractor(e)
		exists, ok := m[k]
		m[k] = resolver(ok, k, exists, v)
	}
	return m
}

// New makes a loop from an abstract source
func New[S, T any](source S, hasNext func(S) bool, getNext func(S) T) func() (T, bool) {
	return func() (out T, ok bool) {
		if hasNext(source) {
			out, ok = getNext(source), true
		}
		return out, ok
	}
}

// RangeClosed creates a loop that generates integers in the range defined by from and to inclusive
func RangeClosed[T constraints.Integer](from T, toInclusive T) func() (T, bool) {
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
func Range[T constraints.Integer](from T, toExclusive T) func() (T, bool) {
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
func OfIndexed[T any](len int, next func(int) T) func() (T, bool) {
	i := 0
	return func() (out T, ok bool) {
		if ok = i < len; ok {
			out = next(i)
			i++
		}
		return out, ok
	}
}

// ConvertAndReduce converts each elements and merges them into one
func ConvertAndReduce[From, To any](next func() (From, bool), converter func(From) To, merger func(To, To) To) (out To) {
	if v, ok := next(); ok {
		out = converter(v)
	} else {
		return out
	}
	for v, ok := next(); ok; v, ok = next() {
		out = merger(out, converter(v))
	}
	return out
}

// ConvAndReduce converts each elements and merges them into one
func ConvAndReduce[From, To any](next func() (From, bool), converter func(From) (To, error), merger func(To, To) To) (out To, err error) {
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
		out = merger(out, c)
	}
	return out, nil
}

func Start[T any](next func() (T, bool)) (func() (T, bool), T, bool) {
	element, ok := next()
	return next, element, ok
}
