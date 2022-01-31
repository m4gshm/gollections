package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func NewSet[T comparable](elements []T) *Set[T] {
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
	return WrapSet(order, uniques)
}

func WrapSet[T comparable](elements []*T, uniques map[T]struct{}) *Set[T] {
	return &Set[T]{elements: elements, uniques: uniques}
}

type Set[T comparable] struct {
	elements []*T
	uniques  map[T]struct{}
}

var _ typ.Set[any] = (*Set[any])(nil)
var _ fmt.Stringer = (*Set[any])(nil)
var _ fmt.GoStringer = (*Set[any])(nil)

func (s *Set[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Set[T]) Iter() *it.RefIter[T] {
	return it.NewRef(&s.elements)
}

func (s *Set[T]) Collect() []T {
	e := s.elements
	out := make([]T, len(e))
	for i, v := range e {
		out[i] = *v
	}
	return out
}

func (s *Set[T]) For(walker func(T) error) error {
	return slice.ForRefs(s.elements, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	slice.ForEachRef(s.elements, walker)
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
	return len(s.elements)
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *Set[T]) String() string {
	return s.GoString()
}

func (s *Set[T]) GoString() string {
	return slice.ToStringRefs(s.elements)
}
