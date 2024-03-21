package loop

import "github.com/m4gshm/gollections/break/loop"

// Loop is a function that returns the next key, value or false if there are no more elements.
type Loop[K, V any] func() (K, V, bool, error)

// Track applies the 'tracker' function to position/element pairs retrieved by the 'next' function. Return the c.ErrBreak to stop tracking..
func (next Loop[K, V]) Track(tracker func(K, V) error) error {
	return loop.Track(next, tracker)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func (next Loop[K, V]) First(predicate func(K, V) bool) (K, V, bool, error) {
	return First(next, predicate)
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merger' function
func (next Loop[K, V]) Reduce(merger func(K, K, V, V) (K, V)) (K, V, error) {
	return Reduce(next, merger)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[K, V]) HasAny(predicate func(K, V) bool) (bool, error) {
	return HasAny(next, predicate)
}

// Filt creates an iterator that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filt(filter func(K, V) (bool, error))  Loop[K, V] {
	return Filt(next, filter)
}

// Filter creates an iterator that checks elements by the 'filter' function and returns successful ones.
func (next Loop[K, V]) Filter(filter func(K, V) bool) Loop[K, V] {
	return Filter(next, filter)
}
