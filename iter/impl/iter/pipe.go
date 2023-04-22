package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// NewPipe instantiates Pipe based on iterator elements.
func NewPipe[T any, IT c.Iterator[T]](iterator IT) Pipe[T, IT] {
	return Pipe[T, IT]{iter: iterator, next: iterator.Next}
}

// Pipe is the Iterator based pipe implementation.
type Pipe[T any, IT any] struct {
	iter IT
	next func() (T, bool)
}

var (
	_ c.Pipe[any]     = (*Pipe[any, c.Iterator[any]])(nil)
	_ c.Iterator[any] = (*Pipe[any, c.Iterator[any]])(nil)
)

// Next implements c.Iterator
func (p Pipe[T, IT]) Next() (element T, ok bool) {
	if next := p.next; next != nil {
		element, ok = next()
	}
	return element, ok
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (p Pipe[T, IT]) Filter(predicate func(T) bool) c.Pipe[T] {
	f := Filter(p.next, predicate)
	return NewPipe[T, c.Iterator[T]](f)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (p Pipe[T, IT]) Convert(converter func(T) T) c.Pipe[T] {
	conv := Convert(p.iter, p.next, converter)
	return NewPipe[T, c.Iterator[T]](conv)
}

// ForEach applies the 'walker' function for every element
func (p Pipe[T, IT]) ForEach(walker func(T)) {
	loop.ForEach(p.next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (p Pipe[T, IT]) For(walker func(T) error) error {
	return loop.For(p.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (p Pipe[T, IT]) Reduce(by func(T, T) T) T {
	return loop.Reduce(p.next, by)
}

// Begin creates iterator
func (p Pipe[T, IT]) Begin() c.Iterator[T] {
	return p
}

// Slice collects the elements to a slice
func (p Pipe[T, IT]) Slice() []T {
	return loop.ToSlice(p.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (p Pipe[T, IT]) HasAny(predicate func(T) bool) bool {
	return loop.HasAny(p.next, predicate)
}
