package iter

import (
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/typ"
)

func NewStream[T any](it typ.Iterator[T]) *IterStream[T] {
	return &IterStream[T]{it: it}
}

type IterStream[T any] struct {
	it typ.Iterator[T]
}

var _ typ.Stream[any] = (*IterStream[any])(nil)

func (s *IterStream[T]) Filter(fit typ.Predicate[T]) typ.Stream[T] {
	return NewStream[T](Filter(s.it, fit))
}

func (s *IterStream[T]) Map(by typ.Converter[T, T]) typ.Stream[T] {
	return NewStream[T](Map(s.it, by))
}

func (s *IterStream[T]) ForEach(w typ.Walker[T]) {
	for s.it.HasNext() {
		w(s.it.Get())
	}
}

func (s *IterStream[T]) Reduce(by op.Binary[T]) T {
	return Reduce(s.it, by)
}

func (s *IterStream[T]) Begin() typ.Iterator[T] {
	return s.it
}
