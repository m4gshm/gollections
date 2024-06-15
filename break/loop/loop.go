package loop

// Loop is a function that returns the next element, false if there are no more elements or error if something is wrong.
type Loop[T any] func() (T, bool, error)

// All is used to iterate through the loop using `for ... range`. Supported since go 1.22 with GOEXPERIMENT=rangefunc enabled.
func (next Loop[T]) All(consumer func(T, error) bool) {
	All(next, consumer)
}

// For applies the 'consumer' function for the elements retrieved by the 'next' function until the consumer returns the c.Break to stop.
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

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func (next Loop[T]) Reduce(merge func(T, T) T) (T, error) {
	return Reduce(next, merge)
}

// ReduceOK reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func (next Loop[T]) ReduceOK(merge func(T, T) T) (T, bool, error) {
	return ReduceOK(next, merge)
}

// Reducee reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func (next Loop[T]) Reducee(merge func(T, T) (T, error)) (T, error) {
	return Reducee(next, merge)
}

// ReduceeOK reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func (next Loop[T]) ReduceeOK(merge func(T, T) (T, error)) (T, bool, error) {
	return ReduceeOK(next, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func (next Loop[T]) Accum(first T, merge func(T, T) T) (T, error) {
	return Accum(first, next, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func (next Loop[T]) Accumm(first T, merge func(T, T) (T, error)) (T, error) {
	return Accumm(first, next, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[T]) HasAny(predicate func(T) bool) (bool, error) {
	return HasAny(next, predicate)
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func (next Loop[T]) Filter(filter func(T) bool) Loop[T] {
	return Filter(next, filter)
}

// Crank rertieves a next element from the 'next' function, returns the function, element, successfully flag.
func (next Loop[T]) Crank() (Loop[T], T, bool, error) {
	return Crank(next)
}
