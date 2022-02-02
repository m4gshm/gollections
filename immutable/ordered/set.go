package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewSet[T comparable](elements []T) *Set[T] {
	var (
		uniques = map[T]struct{}{}
		order   = []T{}
	)
	for _, v := range elements {
		if _, ok := uniques[v]; !ok {
			order = append(order, v)
			uniques[v] = struct{}{}
		}
	}
	return WrapSet(order, uniques)
}

func WrapSet[T comparable](elements []T, uniques map[T]struct{}) *Set[T] {
	return &Set[T]{elements: elements, uniques: uniques}
}

type Set[T comparable] struct {
	elements []T
	uniques  map[T]struct{}
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer   = (*Set[int])(nil)
)

func (s *Set[T]) Begin() c.Iterator[T] {
	return s.Iter()
}

func (s *Set[T]) Iter() *it.PIter[T] {
	return it.NewP(&s.elements)
}

func (s *Set[T]) Collect() []T {
	return slice.Copy(s.elements)
}

func (s *Set[T]) For(walker func(T) error) error {
	return slice.For(s.elements, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	slice.ForEach(s.elements, walker)
}

func (s *Set[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T, c.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Set[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T, c.Iterator[T]] {
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
	return slice.ToString(s.elements)
}
