package loop

import (
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
)

// Loop is a function that returns the next element or ok=false if there are no more elements.
type Loop[T any] func() (element T, ok bool)

var (
	_ c.Filterable[any, Loop[any], loop.Loop[any]]  = (Loop[any])(nil)
	_ c.Convertable[any, Loop[any], loop.Loop[any]] = (Loop[any])(nil)
)

// All is used to iterate through the loop using `for ... range`.
func (next Loop[T]) All(consumer func(T) bool) {
	All(next, consumer)
}

// For applies the 'consumer' function for the elements retrieved by the 'next' function until the consumer returns the c.Break to stop.
func (next Loop[T]) For(consumer func(T) error) error {
	return For(next, consumer)
}

// ForEach applies the 'consumer' function to the elements retrieved by the 'next' function
func (next Loop[T]) ForEach(consumer func(T)) {
	ForEach(next, consumer)
}

// ForEachFiltered applies the 'consumer' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func (next Loop[T]) ForEachFiltered(predicate func(T) bool, consumer func(T)) {
	ForEachFiltered(next, predicate, consumer)
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
func (next Loop[T]) SliceCap(capacity int) []T {
	return SliceCap(next, capacity)
}

// Append collects the elements retrieved by the 'next' function into the specified 'out' slice
func (next Loop[T]) Append(out []T) []T {
	return Append(next, out)
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func (next Loop[T]) Reduce(merge func(T, T) T) T {
	return Reduce(next, merge)
}

// ReduceOK reduces the elements retrieved by the 'next' function into an one using the 'merge' function.
// Returns ok==false if the 'next' function returns ok=false at the first call (no more elements).
func (next Loop[T]) ReduceOK(merge func(T, T) T) (result T, ok bool) {
	return ReduceOK(next, merge)
}

// Reducee reduces the elements retrieved by the 'next' function into an one pair using the 'merge' function.
// If the 'next' function returns ok=false at the first call, the zero value of 'T' type is returned.
func (next Loop[T]) Reducee(merge func(T, T) (T, error)) (T, error) {
	return Reducee(next, merge)
}

// ReduceeOK reduces the elements retrieved by the 'next' function into an one pair using the 'merge' function.
func (next Loop[T]) ReduceeOK(merge func(T, T) (T, error)) (resul T, ok bool, err error) {
	return ReduceeOK(next, merge)
}

// Accum accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func (next Loop[T]) Accum(first T, merge func(T, T) T) T {
	return Accum(first, next, merge)
}

// Accumm accumulates a value by using the 'first' argument to initialize the accumulator and sequentially applying the 'merge' functon to the accumulator and each element retrieved by the 'next' function.
func (next Loop[T]) Accumm(first T, merge func(T, T) (T, error)) (T, error) {
	return Accumm(first, next, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (next Loop[T]) HasAny(predicate func(T) bool) bool {
	return HasAny(next, predicate)
}

// Filt creates a loop that checks elements by the 'filter' function and returns successful ones.
func (next Loop[T]) Filt(filter func(T) (bool, error)) loop.Loop[T] {
	return Filt(next, filter)
}

// Filter creates a loop that checks elements by the 'filter' function and returns successful ones.
func (next Loop[T]) Filter(filter func(T) bool) Loop[T] {
	return Filter(next, filter)
}

// Convert creates a loop that applies the 'converter' function to iterable elements.
func (next Loop[T]) Convert(converter func(T) T) Loop[T] {
	return Convert(next, converter)
}

// Conv creates a loop that applies the 'converter' function to iterable elements.
func (next Loop[T]) Conv(converter func(T) (T, error)) loop.Loop[T] {
	return Conv(next, converter)
}

// Crank rertieves a next element from the 'next' function, returns the function, element, successfully flag.
func (next Loop[T]) Crank() (Loop[T], T, bool) {
	return Crank(next)
}
