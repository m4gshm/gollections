package loop

// Loop is a function that returns the next element, false if there are no more elements or error if something is wrong.
type Loop[T any] func() (T, bool, error)

// For applies the 'consumer' function for the elements retrieved by the 'next' function. Return the c.Break to stop
func (next Loop[T]) For(consumer func(T) error) error {
	return For(next, consumer)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func (next Loop[T]) First(predicate func(T) bool) (T, bool, error) {
	return First(next, predicate)
}

// Slice collects the elements retrieved by the 'next' function into a new slice
func (next Loop[T]) Slice() ([]T, error) {
	return Slice(next)
}

// SliceCap collects the elements retrieved by the 'next' function into a new slice with predefined capacity
func (next Loop[T]) SliceCap(cap int) ([]T, error) {
	return SliceCap(next, cap)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func (next Loop[T]) Append(out []T) ([]T, error) {
	return Append(next, out)
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merger' function
func (next Loop[T]) Reduce(merger func(T, T) T) (T, error) {
	return Reduce(next, merger)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[T]) HasAny(predicate func(T) bool) (bool, error) {
	return HasAny(next, predicate)
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func (next Loop[T]) Filter(filter func(T) bool) Loop[T] {
	return Filter(next, filter)
}
