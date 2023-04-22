package transform

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
)

// New instantiates Pipe based on iterator elements.
func New[T any](next func() (T, bool)) Transform[T] {
	return Transform[T]{next: next}
}

// Transform is the Iterator based pipe implementation.
type Transform[T any] struct {
	next func() (T, bool)
}

var (
	_ c.Transform[any] = (*Transform[any])(nil)
	_ c.Transform[any] = Transform[any]{}
)

// Next implements c.Iterator
func (t Transform[T]) Next() (element T, ok bool) {
	if next := t.next; next != nil {
		element, ok = next()
	}
	return element, ok
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (t Transform[T]) Filter(predicate func(T) bool) c.Transform[T] {
	f := iter.Filter(t.next, predicate)
	return New(f.Next)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (t Transform[T]) Convert(converter func(T) T) c.Transform[T] {
	conv := iter.Convert(t.next, converter)
	return New(conv.Next)
}

// ForEach applies the 'walker' function for every element
func (t Transform[T]) ForEach(walker func(T)) {
	loop.ForEach(t.next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (t Transform[T]) For(walker func(T) error) error {
	return loop.For(t.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (t Transform[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(t.next, by)
}

// Begin creates iterator
func (t Transform[T]) Begin() c.Iterator[T] {
	return t
}

// Slice collects the elements to a slice
func (t Transform[T]) Slice() []T {
	return loop.ToSlice(t.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (t Transform[T]) HasAny(predicate func(T) bool) bool {
	return loop.HasAny(t.next, predicate)
}
