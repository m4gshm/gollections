package it

import (
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func NewPipe[T any](it typ.Iterator[T]) *IterPipe[T] {
	return &IterPipe[T]{it: it}
}

type IterPipe[T any] struct {
	it       typ.Iterator[T]
	elements []T
}

var _ typ.Pipe[any, []any, typ.Iterator[any]] = (*IterPipe[any])(nil)

func (s *IterPipe[T]) Filter(fit typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return NewPipe[T](Filter(s.it, fit))
}

func (s *IterPipe[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
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

func (s *IterPipe[T]) Begin() typ.Iterator[T] {
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
