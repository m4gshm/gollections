package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// NewPipe instantiates Pipe based on iterator elements.
func NewPipe[T any, IT c.Iterator[T]](iterator IT) *Pipe[T] {
	return &Pipe[T]{Iterator: iterator}
}

// Pipe is the Iterator based pipe implementation.
type Pipe[T any] struct {
	c.Iterator[T]
}

var _ c.Pipe[any] = (*Pipe[any])(nil)

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (s *Pipe[T]) Filter(predicate func(T) bool) c.Pipe[T] {
	return NewPipe[T](Filter(s, s.Next, predicate))
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (s *Pipe[T]) Convert(converter func(T) T) c.Pipe[T] {
	return NewPipe[T](Convert(s, s.Next, converter))
}

// ForEach applies the 'walker' function for every element
func (s *Pipe[T]) ForEach(walker func(T)) {
	loop.ForEach(s.Next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (s *Pipe[T]) For(walker func(T) error) error {
	return loop.For(s.Next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s *Pipe[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(s.Next, by)
}

// Begin creates iterator
func (s *Pipe[T]) Begin() c.Iterator[T] {
	return s
}

// Slice collects the elements to a slice
func (s *Pipe[T]) Slice() []T {
	return loop.ToSlice(s.Next)
}
