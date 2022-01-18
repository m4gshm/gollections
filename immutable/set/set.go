package set

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func Convert[T comparable](elements []T) *Set[T] {
	uniques := make(map[T]struct{}, 0)
	for _, v := range elements {
		uniques[v] = struct{}{}
	}
	return Wrap(uniques)
}

func Wrap[T comparable](uniques map[T]struct{}) *Set[T] {
	return &Set[T]{uniques: uniques}
}

type Set[T comparable] struct {
	uniques map[T]struct{}
}

var _ typ.Set[any, typ.Iterator[any]] = (*Set[any])(nil)
var _ fmt.Stringer = (*Set[any])(nil)
var _ fmt.GoStringer = (*Set[any])(nil)

func (s *Set[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Set[T]) Iter() *SetIter[T] {
	return &SetIter[T]{KV: it.NewKV(s.uniques)}
}

func (s *Set[T]) Collect() []T {
	e := s.uniques
	out := make([]T, 0, len(e))
	for k := range e {
		out = append(out, k)
	}
	return out
}

func (s *Set[T]) For(walker func(T) error) error {
	for k := range s.uniques {
		if err := walker(k); err != nil {
			return err
		}
	}
	return nil
}

func (s *Set[T]) ForEach(walker func(T)) error {
	return s.For(func(t T) error { walker(t); return nil })
}

func (s *Set[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Set[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *Set[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *Set[T]) Len() int {
	return len(s.uniques)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *Set[T]) String() string {
	return s.GoString()
}

func (s *Set[T]) GoString() string {
	return slice.ToString(s.Collect())
}
