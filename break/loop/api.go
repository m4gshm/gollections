// Package loop provides helpers for loop operation and iterator implementations
package loop

import (
	"errors"

	"github.com/m4gshm/gollections/break/op"
	"github.com/m4gshm/gollections/break/predicate/always"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/convert"
	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/map_/resolv"
	"github.com/m4gshm/gollections/notsafe"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = c.ErrBreak

// From wrap the next loop to a breakable loop
func From[T any](next func() (T, bool)) func() (T, bool, error) {
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
func ForFiltered[T any](next func() (T, bool, error), walker func(T) error, predicate func(T) (bool, error)) error {
	for {
		if v, ok, err := next(); err != nil || !ok {
			return err
		} else if ok, err := predicate(v); err != nil {
			return err
		} else if ok {
			if err := walker(v); err != nil {
				return brk(err)
			}
		}
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool, error), predicate func(T) (bool, error)) (T, bool, error) {
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
	for {
		v, ok, err := next()
		if err != nil || !ok {
			return out, err
		}
		out = append(out, v)
	}
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function
func Reduce[T any](next func() (T, bool, error), merger func(T, T) (T, error)) (out T, e error) {
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
func HasAny[T any](next func() (T, bool, error), predicate func(T) (bool, error)) (bool, error) {
	_, ok, err := First(next, predicate)
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
func Conv[From, To any](next func() (From, bool, error), converter func(From) (To, error)) ConvertIter[From, To] {
	return ConvertIter[From, To]{next: next, converter: converter}
}

// Convert instantiates an iterator that converts elements with a converter and returns them.
func Convert[From, To any](next func() (From, bool, error), converter func(From) To) ConvertIter[From, To] {
	return ConvertIter[From, To]{next: next, converter: func(f From) (To, error) { return converter(f), nil }}
}

// ConvCheck is similar to ConvertFit, but it checks and transforms elements together
func ConvCheck[From, To any](next func() (From, bool, error), converter func(from From) (To, bool, error)) ConvertCheckIter[From, To] {
	return ConvertCheckIter[From, To]{next: next, converter: converter}
}

// ConvertCheck is similar to ConvFit, but it checks and transforms elements together
func ConvertCheck[From, To any](next func() (From, bool, error), converter func(from From) (To, bool)) ConvertCheckIter[From, To] {
	return ConvertCheckIter[From, To]{next: next, converter: func(f From) (To, bool, error) { c, ok := converter(f); return c, ok, nil }}
}

// FitAndConv returns a stream that filters source elements and converts them
func FitAndConv[From, To any](next func() (From, bool, error), filter func(From) (bool, error), converter func(From) (To, error)) ConvertFitIter[From, To] {
	return FilterConvertFilter(next, filter, converter, always.True[To])
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[From, To any](next func() (From, bool, error), filter func(From) bool, converter func(From) To) ConvertFitIter[From, To] {
	return FilterConvertFilter(next, func(f From) (bool, error) { return filter(f), nil }, func(f From) (To, error) { return converter(f), nil }, always.True[To])
}

// FilterConvertFilter filters source, converts, and filters converted elements
func FilterConvertFilter[From, To any](next func() (From, bool, error), filter func(From) (bool, error), converter func(From) (To, error), filterTo func(To) (bool, error)) ConvertFitIter[From, To] {
	return ConvertFitIter[From, To]{next: next, converter: converter, filterFrom: filter, filterTo: filterTo}
}

// ConvertAndFilter additionally filters 'To' elements
func ConvertAndFilter[From, To any](next func() (From, bool, error), converter func(From) (To, error), filter func(To) (bool, error)) ConvertFitIter[From, To] {
	return FilterConvertFilter(next, always.True[From], converter, filter)
}

// Flat instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flat[From, To any](next func() (From, bool, error), flattener func(From) ([]To, error)) FlatIter[From, To] {
	return FlatIter[From, To]{next: next, flattener: flattener, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Flatt instantiates an iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](next func() (From, bool, error), flattener func(From) []To) FlatIter[From, To] {
	return FlatIter[From, To]{next: next, flattener: func(f From) ([]To, error) { return flattener(f), nil }, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FitAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FitAndFlat[From, To any](next func() (From, bool, error), filter func(From) (bool, error), flattener func(From) ([]To, error)) FlattenFitIter[From, To] {
	return FitFlatFit(next, filter, flattener, always.True[To])
}

// FilterAndFlatt filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlatt[From, To any](next func() (From, bool, error), filter func(From) bool, flattener func(From) []To) FlattenFitIter[From, To] {
	return FitFlatFit(next, func(f From) (bool, error) { return filter(f), nil }, func(f From) ([]To, error) { return flattener(f), nil }, always.True[To])
}

// FlatAndFit extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlatAndFit[From, To any](next func() (From, bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) FlattenFitIter[From, To] {
	return FitFlatFit(next, always.True[From], flattener, filterTo)
}

// FlattAndFilter extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FlattAndFilter[From, To any](next func() (From, bool, error), flattener func(From) []To, filterTo func(To) bool) FlattenFitIter[From, To] {
	return FitFlatFit(next, always.True[From], func(f From) ([]To, error) { return flattener(f), nil }, func(t To) (bool, error) { return filterTo(t), nil })
}

// FitFlatFit filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FitFlatFit[From, To any](next func() (From, bool, error), filterFrom func(From) (bool, error), flattener func(From) ([]To, error), filterTo func(To) (bool, error)) FlattenFitIter[From, To] {
	return FlattenFitIter[From, To]{next: next, filterFrom: filterFrom, flatt: flattener, filterTo: filterTo, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FilterFlattFilter filters source elements, extracts slices of 'To' by the 'flattener' function and filters extracted elements
func FilterFlattFilter[From, To any](next func() (From, bool, error), filterFrom func(From) bool, flattener func(From) []To, filterTo func(To) bool) FlattenFitIter[From, To] {
	return FlattenFitIter[From, To]{
		next:       next,
		filterFrom: func(f From) (bool, error) { return filterFrom(f), nil },
		flatt:      func(f From) ([]To, error) { return flattener(f), nil },
		filterTo:   func(t To) (bool, error) { return filterTo(t), nil },
		elemSizeTo: notsafe.GetTypeSize[To](),
	}
}

// Filt creates an Iterator that checks elements by the 'filter' function and returns successful ones.
func Filt[T any](next func() (T, bool, error), filter func(T) (bool, error)) FiltIter[T] {
	return FiltIter[T]{next: next, filter: filter}
}

// Filter creates an Iterator that checks elements by the 'filter' function and returns successful ones.
func Filter[T any](next func() (T, bool, error), filter func(T) bool) FiltIter[T] {
	return FiltIter[T]{next: next, filter: func(t T) (bool, error) { return filter(t), nil }}
}

// NotNil creates an Iterator that filters nullable elements.
func NotNil[T any](next func() (*T, bool, error)) FiltIter[*T] {
	return Filt(next, as.ErrTail(check.NotNil[T]))
}

// ToKV transforms iterable elements to key/value iterator based on applying key, value extractors to the elements
func ToKV[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) KeyValuer[T, K, V] {
	kv := NewKeyValuer(next, keyExtractor, valExtractor)
	return kv
}

// ToKVs transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func ToKVs[T, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valsExtractor func(T) ([]V, error)) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(next, keysExtractor, valsExtractor)
	return &kv
}

// FlattKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattKeys[T, K, V any](next func() (T, bool, error), keysExtractor func(T) []K, valExtractor func(T) V) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(next, func(t T) ([]K, error) { return keysExtractor(t), nil }, convSlice(func(t T) (V, error) { return valExtractor(t), nil }))
	return &kv
}

// FlatKeys transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlatKeys[T, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valExtractor func(T) (V, error)) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(next, keysExtractor, convSlice(valExtractor))
	return &kv
}

// FlattValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlattValues[T, K, V any](next func() (T, bool, error), keyExtractor func(T) K, valsExtractor func(T) []V) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(next, convSlice(func(t T) (K, error) { return keyExtractor(t), nil }), func(t T) ([]V, error) { return valsExtractor(t), nil })
	return &kv
}

// FlatValues transforms iterable elements to key/value iterator based on applying key, value extractor to the elements
func FlatValues[T, K, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valsExtractor func(T) ([]V, error)) *MultipleKeyValuer[T, K, V] {
	kv := NewMultipleKeyValuer(next, convSlice(keyExtractor), valsExtractor)
	return &kv
}

// Group converts elements retrieved by the 'next' function into a map, extracting a key for each element applying the converter 'keyExtractor'.
// The keyExtractor converts an element to a key.
// The valExtractor converts an element to an value.
func Group[T any, K comparable, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valExtractor func(T) (V, error)) (map[K][]V, error) {
	return ToMapResolv(next, keyExtractor, valExtractor, func(ok bool, k K, rv []V, v V) ([]V, error) {
		return resolv.Append(ok, k, rv, v), nil
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

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[T any, K comparable, V, VR any](
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

// New is the main breakable loop constructor
func New[S, T any](source S, hasNext func(S) bool, getNext func(S) (T, error)) func() (T, bool, error) {
	return func() (out T, ok bool, err error) {
		if ok := hasNext(source); !ok {
			return out, false, nil
		} else if n, err := getNext(source); err != nil {
			return out, false, err
		} else {
			return n, true, nil
		}
	}
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
