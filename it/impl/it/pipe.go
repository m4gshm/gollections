package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// NewPipe instantiates Pipe based on iterator elements.
func NewPipe[T any, IT c.Iterator[T]](iter IT) *IterPipe[T] {
	return &IterPipe[T]{Iterator: iter}
}

// IterPipe is the Iterator based pipe implementation.
type IterPipe[T any] struct {
	c.Iterator[T]
}

var _ c.Pipe[any, []any] = (*IterPipe[any])(nil)

func (s *IterPipe[T]) Filter(filter func(T) bool) c.Pipe[T, []T] {
	it := s
	return NewPipe[T](Filter(it, it.Next, filter))
}

func (s *IterPipe[T]) Convert(by func(T) T) c.Pipe[T, []T] {
	return NewPipe[T](Convert(s, s.Next, by))
}

func (s *IterPipe[T]) ForEach(walker func(T)) {
	loop.ForEach(s.Next, walker)
}

func (s *IterPipe[T]) For(walker func(T) error) error {
	return loop.For(s.Next, walker)
}

func (s *IterPipe[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(s.Next, by)
}

func (s *IterPipe[T]) Begin() c.Iterator[T] {
	return s
}

func (s *IterPipe[T]) Collect() []T {
	var e []T
	for v, ok := s.Next(); ok; v, ok = s.Next() {
		e = append(e, v)
	}
	return e
}
