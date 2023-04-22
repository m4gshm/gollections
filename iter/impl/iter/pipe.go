package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// NewPipe instantiates Pipe based on iterator elements.
func NewPipe[T any](next func() (T, bool)) Pipe[T] {
	return Pipe[T]{next: next}
}

// Pipe is the Iterator based pipe implementation.
type Pipe[T any] struct {
	next func() (T, bool)
}

var (
	_ c.Pipe[any]     = (*Pipe[any])(nil)
	_ c.Iterator[any] = (*Pipe[any])(nil)
	_ c.Pipe[any]     = Pipe[any]{}
	_ c.Iterator[any] = Pipe[any]{}
)

// Next implements c.Iterator
func (p Pipe[T]) Next() (element T, ok bool) {
	if next := p.next; next != nil {
		element, ok = next()
	}
	return element, ok
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (p Pipe[T]) Filter(predicate func(T) bool) c.Pipe[T] {
	f := Filter(p.next, predicate)
	return NewPipe(f.Next)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (p Pipe[T]) Convert(converter func(T) T) c.Pipe[T] {
	conv := Convert(p.next, converter)
	return NewPipe(conv.Next)
}

// ForEach applies the 'walker' function for every element
func (p Pipe[T]) ForEach(walker func(T)) {
	loop.ForEach(p.next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (p Pipe[T]) For(walker func(T) error) error {
	return loop.For(p.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (p Pipe[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(p.next, by)
}

// Begin creates iterator
func (p Pipe[T]) Begin() c.Iterator[T] {
	return p
}

// Slice collects the elements to a slice
func (p Pipe[T]) Slice() []T {
	return loop.ToSlice(p.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (p Pipe[T]) HasAny(predicate func(T) bool) bool {
	return loop.HasAny(p.next, predicate)
}
