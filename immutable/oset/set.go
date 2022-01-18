package oset

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func ToOrderedSet[T comparable](elements []T) *OrderedSet[T] {
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
	return WrapOrderedSet(order, uniques)
}

func NewOrderedSet[T comparable](capacity int) *OrderedSet[T] {
	return WrapOrderedSet(make([]*T, 0, capacity), make(map[T]struct{}, capacity))
}

func WrapOrderedSet[T comparable](elements []*T, uniques map[T]struct{}) *OrderedSet[T] {
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
	return &it.RefIter[T]{Iterator: s.newIter()}
}

func (s *OrderedSet[T]) newIter() *it.Iter[*T] {
	return it.New(s.elements)
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
	for _, e := range s.elements {
		if err := walker(*e); err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderedSet[T]) ForEach(walker func(T)) error {
	return s.For(func(t T) error { walker(t); return nil })
}

func (s *OrderedSet[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](&it.RefIter[T]{Iterator: it.Filter(s.newIter(), func(ref *T) bool { return filter(*ref) })})
}

func (s *OrderedSet[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](&it.RefIter[T]{Iterator: it.Map(s.newIter(), func(ref *T) *T {
		conv := by(*ref)
		return &conv
	})})
}

func (s *OrderedSet[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(&it.RefIter[T]{Iterator: s.newIter()}, by)
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
