package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/predicate"
)

// NewPipe instantiates Pipe based on iterable elements.
func NewPipe[T any, IT c.Iterator[T]](iter IT) *IterPipe[T] {
	return &IterPipe[T]{it: iter}
}

// IterPipe is the Iterator based pipe implementation.
type IterPipe[T any] struct {
	it       c.Iterator[T]
	elements []T
}

var _ c.Pipe[any, []any] = (*IterPipe[any])(nil)

func (s *IterPipe[T]) Filter(fit predicate.Predicate[T]) c.Pipe[T, []T] {
	return NewPipe[T](Filter(s.it, fit))
}

func (s *IterPipe[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	return NewPipe[T](Map(s.it, by))
}

func (s *IterPipe[T]) ForEach(walker func(T)) {
	ForEach(s.it, walker)
}

func (s *IterPipe[T]) For(walker func(T) error) error {
	return For(s.it, walker)
}

func (s *IterPipe[T]) Reduce(by c.Binary[T]) T {
	return Reduce(s.it, by)
}

func (s *IterPipe[T]) Begin() c.Iterator[T] {
	return s.it
}

func (s *IterPipe[T]) Collect() []T {
	e := s.elements
	if e == nil {
		for v, ok := s.it.Next(); ok; v, ok = s.it.Next() {
			e = append(e, v)
		}
		s.elements = e
	}
	return e
}
