package loop

import (
	breakMapConvert "github.com/m4gshm/gollections/break/kv/convert"
	breakkvloop "github.com/m4gshm/gollections/break/kv/loop"
	breakMapFilter "github.com/m4gshm/gollections/break/kv/predicate"
	"github.com/m4gshm/gollections/kv/convert"
	kvPredicate "github.com/m4gshm/gollections/kv/predicate"
)

// Loop is a function that returns the next key\value or ok==false if there are no more elements.
//
// Deprecated: replaced by [github.com/m4gshm/gollections/seq.Seq2]
type Loop[K, V any] func() (key K, value V, ok bool)

// All is used to iterate through the loop using `for ... range`.
func (next Loop[K, V]) All(consumer func(key K, value V) bool) {
	All(next, consumer)
}

// Track applies the 'consumer' function to position/element pairs retrieved by the 'next' function until the consumer returns the c.Break to stop.
func (next Loop[K, V]) Track(consumer func(K, V) error) error {
	return Track(next, consumer)
}

// First returns the first element that satisfies the condition of the 'predicate' function.
func (next Loop[K, V]) First(predicate func(K, V) bool) (K, V, bool) {
	return First(next, predicate)
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
func (next Loop[K, V]) Reduce(merge func(K, K, V, V) (K, V)) (K, V, bool) {
	return ReduceOK(next, merge)
}

// Reducee reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
func (next Loop[K, V]) Reducee(merge func(K, K, V, V) (K, V, error)) (K, V, bool, error) {
	return ReduceeOK(next, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[K, V]) HasAny(predicate func(K, V) bool) bool {
	return HasAny(next, predicate)
}

// Filt creates a seq that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filt(filter func(K, V) (bool, error)) breakkvloop.Loop[K, V] {
	return Filt(next, filter)
}

// Filter creates a seq that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filter(filter func(K, V) bool) Loop[K, V] {
	return Filter(next, filter)
}

// Convert creates a seq that applies the 'converter' function to iterable key\values.
func (next Loop[K, V]) Convert(converter func(K, V) (K, V)) Loop[K, V] {
	return Convert(next, converter)
}

// Conv creates a seq that applies the 'converter' function to iterable key\values.
func (next Loop[K, V]) Conv(converter func(K, V) (K, V, error)) breakkvloop.Loop[K, V] {
	return Conv(next, converter)
}

// FilterKey returns a seq consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FilterKey(predicate func(K) bool) Loop[K, V] {
	return Filter(next, kvPredicate.Key[V](predicate))
}

// FiltKey returns a seq consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FiltKey(predicate func(K) (bool, error)) breakkvloop.Loop[K, V] {
	return Filt(next, breakMapFilter.Key[V](predicate))
}

// ConvertKey returns a seq that applies the 'converter' function to keys of the map
func (next Loop[K, V]) ConvertKey(by func(K) K) Loop[K, V] {
	return Convert(next, convert.Key[V](by))
}

// ConvKey returns a seq that applies the 'converter' function to keys of the map
func (next Loop[K, V]) ConvKey(converter func(K) (K, error)) breakkvloop.Loop[K, V] {
	return Conv(next, breakMapConvert.Key[V](converter))
}

// FilterValue returns a seq consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FilterValue(predicate func(V) bool) Loop[K, V] {
	return Filter(next, kvPredicate.Value[K](predicate))
}

// FiltValue returns a errorable seq consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FiltValue(predicate func(V) (bool, error)) breakkvloop.Loop[K, V] {
	return Filt(next, breakMapFilter.Value[K](predicate))
}

// ConvertValue returns a seq that applies the 'converter' function to values of the map
func (next Loop[K, V]) ConvertValue(converter func(V) V) Loop[K, V] {
	return Convert(next, convert.Value[K](converter))
}

// ConvValue returns a errorable seq that applies the 'converter' function to values of the map
func (next Loop[K, V]) ConvValue(converter func(V) (V, error)) breakkvloop.Loop[K, V] {
	return Conv(next, breakMapConvert.Value[K](converter))
}

// Crank rertieves a next element from the 'next' function, returns the function, element, successfully flag.
func (next Loop[K, V]) Crank() (Loop[K, V], K, V, bool) {
	return Crank(next)
}
