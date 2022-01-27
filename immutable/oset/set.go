package oset

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func Convert[T comparable](elements []T) *OrderedSet[T] {
	var (
		uniques = make(map[T]struct{}, 0)
		order   = make([]*T, 0, 0)
	)
	for _, v := range elements {
		if _, ok := uniques[v]; !ok {
			vv := v
			order = append(order, &vv)
			uniques[vv] = struct{}{}
		}
	}
	return Wrap(order, uniques)
}

func Wrap[T comparable](elements []*T, uniques map[T]struct{}) *OrderedSet[T] {
	return &OrderedSet[T]{elements: elements, uniques: uniques}
}

type OrderedSet[T comparable] struct {
	elements []*T
	uniques  map[T]struct{}
}

var _ typ.Set[any, typ.Iterator[any]] = (*OrderedSet[any])(nil)
var _ fmt.Stringer = (*OrderedSet[any])(nil)
var _ fmt.GoStringer = (*OrderedSet[any])(nil)

func (s *OrderedSet[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *OrderedSet[T]) Iter() *it.RefIter[T] {
	return it.NewRef(&s.elements)
}

func (s *OrderedSet[T]) Collect() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *OrderedSet[T]) For(walker func(T) error) error {
	return slice.ForRefs(s.elements, walker)
}

func (s *OrderedSet[T]) ForEach(walker func(T)) {
	slice.ForEachRef(s.elements, walker)
}

func (s *OrderedSet[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *OrderedSet[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *OrderedSet[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *OrderedSet[T]) Len() int {
	return len(s.elements)
}

func (s *OrderedSet[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *OrderedSet[T]) String() string {
	return s.GoString()
}

func (s *OrderedSet[T]) GoString() string {
	return slice.ToStringRefs(s.elements)
}