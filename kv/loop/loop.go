package loop

import (
	breakkvloop "github.com/m4gshm/gollections/break/kv/loop"
	breakMapFilter "github.com/m4gshm/gollections/break/kv/predicate"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	"github.com/m4gshm/gollections/kv/convert"
	kvPredicate "github.com/m4gshm/gollections/kv/predicate"
)

// Loop is a function that returns the next key, value or false if there are no more elements.
type Loop[K, V any] func() (K, V, bool)

func (next Loop[K, V])  All(consumer func(key K, value V) bool) {
	All(next, consumer)
}

// Track applies the 'tracker' function to position/element pairs retrieved by the 'next' function. Return the c.Break to stop tracking..
func (next Loop[K, V]) Track(tracker func(K, V) error) error {
	return Track(next, tracker)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func (next Loop[K, V]) First(predicate func(K, V) bool) (K, V, bool) {
	return First(next, predicate)
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merger' function
func (next Loop[K, V]) Reduce(merger func(K, K, V, V) (K, V)) (K, V) {
	return Reduce(next, merger)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[K, V]) HasAny(predicate func(K, V) bool) bool {
	return HasAny(next, predicate)
}

// Filt creates a loop that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filt(filter func(K, V) (bool, error)) breakkvloop.Loop[K, V] {
	return Filt(next, filter)
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filter(filter func(K, V) bool) Loop[K, V] {
	return Filter(next, filter)
}

// Convert creates a loop that applies the 'converter' function to iterable key\values.
func (next Loop[K, V]) Convert(converter func(K, V) (K, V)) Loop[K, V] {
	return Convert(next, converter)
}

// Conv creates a loop that applies the 'converter' function to iterable key\values.
func (next Loop[K, V]) Conv(converter func(K, V) (K, V, error)) breakkvloop.Loop[K, V] {
	return Conv(next, converter)
}

// FilterKey returns a loop consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FilterKey(predicate func(K) bool) Loop[K, V] {
	return Filter(next, kvPredicate.Key[V](predicate))
}

// FiltKey returns a loop consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FiltKey(predicate func(K) (bool, error)) breakkvloop.Loop[K, V] {
	return Filt(next, breakMapFilter.Key[V](predicate))
}

// ConvertKey returns a loop that applies the 'converter' function to keys of the map
func (next Loop[K, V]) ConvertKey(by func(K) K) Loop[K, V] {
	return Convert(next, convert.Key[V](by))
}

// ConvKey returns a loop that applies the 'converter' function to keys of the map
func (next Loop[K, V]) ConvKey(converter func(K) (K, error)) breakkvloop.Loop[K, V] {
	return Conv(next, breakMapConvert.Key[V](converter))
}

// FilterValue returns a loop consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FilterValue(predicate func(V) bool) Loop[K, V] {
	return Filter(next, kvPredicate.Value[K](predicate))
}

// FiltValue returns a breakable loop consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (next Loop[K, V]) FiltValue(predicate func(V) (bool, error)) breakkvloop.Loop[K, V] {
	return Filt(next, breakMapFilter.Value[K](predicate))
}

// ConvertValue returns a loop that applies the 'converter' function to values of the map
func (next Loop[K, V]) ConvertValue(converter func(V) V) Loop[K, V] {
	return Convert(next, convert.Value[K](converter))
}

// ConvValue returns a breakable loop that applies the 'converter' function to values of the map
func (next Loop[K, V]) ConvValue(converter func(V) (V, error)) breakkvloop.Loop[K, V] {
	return Conv(next, breakMapConvert.Value[K](converter))
}
