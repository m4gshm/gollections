package loop

import "github.com/m4gshm/gollections/break/loop"

// Loop is a function that returns the next element or false if there are no more elements.
type Loop[T any] func() (T, bool)

// For applies the 'walker' function for the elements retrieved by the 'next' function. Return the c.ErrBreak to stop
func (next Loop[T]) For(walker func(T) error) error {
	return For(next, walker)
}

// ForEach applies the 'walker' function to the elements retrieved by the 'next' function
func (next Loop[T]) ForEach(walker func(T)) {
	ForEach(next, walker)
}

// ForEachFiltered applies the 'walker' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func (next Loop[T]) ForEachFiltered(predicate func(T) bool, walker func(T)) {
	ForEachFiltered(next, predicate, walker)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func (next Loop[T]) First(predicate func(T) bool) (T, bool) {
	return First(next, predicate)
}

// Slice collects the elements retrieved by the 'next' function into a new slice
func (next Loop[T]) Slice() []T {
	return Slice(next)
}

// SliceCap collects the elements retrieved by the 'next' function into a new slice with predefined capacity
func (next Loop[T]) SliceCap(cap int) []T {
	return SliceCap(next, cap)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func (next Loop[T]) Append(out []T) []T {
	return Append(next, out)
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merger' function
func (next Loop[T]) Reduce(merger func(T, T) T) T {
	return Reduce(next, merger)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[T]) HasAny(predicate func(T) bool) bool {
	return HasAny(next, predicate)
}

// Filt creates an iterator that checks elements by the 'filter' function and returns successful ones.
func (next Loop[T]) Filt(filter func(T) (bool, error)) loop.Loop[T] {
	return Filt(next, filter)
}

// Filter creates an iterator that checks elements by the 'filter' function and returns successful ones.
func (next Loop[T]) Filter(filter func(T) bool) Loop[T] {
	return Filter(next, filter)
}
