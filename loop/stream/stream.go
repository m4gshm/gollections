package stream

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/loop/iter"
)

// New instantiates Pipe based on iterator elements.
func New[T any](next func() (T, bool)) Stream[T] {
	return Stream[T]{next: next}
}

// Stream is the Iterator based pipe implementation.
type Stream[T any] struct {
	next func() (T, bool)
}

var (
	_ c.Stream[any] = (*Stream[any])(nil)
	_ c.Stream[any] = Stream[any]{}
)

// Next implements c.Iterator
func (t Stream[T]) Next() (element T, ok bool) {
	if next := t.next; next != nil {
		element, ok = next()
	}
	return element, ok
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (t Stream[T]) Filter(predicate func(T) bool) c.Stream[T] {
	f := iter.Filter(t.next, predicate)
	return New(f.Next)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (t Stream[T]) Convert(converter func(T) T) c.Stream[T] {
	conv := iter.Convert(t.next, converter)
	return New(conv.Next)
}

// ForEach applies the 'walker' function for every element
func (t Stream[T]) ForEach(walker func(T)) {
	loop.ForEach(t.next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (t Stream[T]) For(walker func(T) error) error {
	return loop.For(t.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (t Stream[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(t.next, by)
}

// Begin creates iterator
func (t Stream[T]) Begin() c.Iterator[T] {
	return t
}

// Slice collects the elements to a slice
func (t Stream[T]) Slice() []T {
	return loop.ToSlice(t.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (t Stream[T]) HasAny(predicate func(T) bool) bool {
	return loop.HasAny(t.next, predicate)
}
