package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewPipe[T any](it c.Iterator[T]) *IterPipe[T] {
	return &IterPipe[T]{it: it}
}

type IterPipe[T any] struct {
	it       c.Iterator[T]
	elements []T
}

var _ c.Pipe[any, []any, c.Iterator[any]] = (*IterPipe[any])(nil)

func (s *IterPipe[T]) Filter(fit c.Predicate[T]) c.Pipe[T, []T, c.Iterator[T]] {
	return NewPipe[T](Filter(s.it, fit))
}

func (s *IterPipe[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T, c.Iterator[T]] {
	return NewPipe[T](Map(s.it, by))
}

func (s *IterPipe[T]) For(walker func(T) error) error {
	for s.it.HasNext() {
		if n, err := s.it.Get(); err != nil {
			return err
		} else if err = walker(n); err != nil {
			return err
		}
	}
	return nil
}

func (s *IterPipe[T]) ForEach(walker func(T)) {
	slice.ForEach(s.elements, walker)
}

func (s *IterPipe[T]) Reduce(by op.Binary[T]) T {
	return Reduce(s.it, by)
}

func (s *IterPipe[T]) Begin() c.Iterator[T] {
	return s.it
}

func (s *IterPipe[T]) Collect() []T {
	e := s.elements
	if e == nil {
		e = make([]T, 0)
		it := s.it
		for it.HasNext() {
			n, err := it.Get()
			if err != nil {
				panic(err)
			}
			e = append(e, n)
		}
		s.elements = e
	}
	return e
}
