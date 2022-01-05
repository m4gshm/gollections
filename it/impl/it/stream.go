package it

import (
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/typ"
)

func NewPipe[T any](it typ.Iterator[T]) *IterPipe[T] {
	return &IterPipe[T]{it: it}
}

type IterPipe[T any] struct {
	it typ.Iterator[T]
}

var _ typ.Pipe[any, typ.Iterator[any]] = (*IterPipe[any])(nil)

func (s *IterPipe[T]) Filter(fit typ.Predicate[T]) typ.Pipe[T, typ.Iterator[T]] {
	return NewPipe[T](Filter(s.it, fit))
}

func (s *IterPipe[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, typ.Iterator[T]] {
	return NewPipe[T](Map(s.it, by))
}

func (s *IterPipe[T]) ForEach(walker func(T)) {
	for s.it.HasNext() {
		walker(s.it.Get())
	}
}

func (s *IterPipe[T]) Reduce(by op.Binary[T]) T {
	return Reduce(s.it, by)
}

func (s *IterPipe[T]) Begin() typ.Iterator[T] {
	return s.it
}
