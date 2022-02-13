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
	return &Set[T]{elements: elements, uniques: uniques, esize: it.GetTypeSize[T]()}
}

type Set[T comparable] struct {
	elements []T
	uniques  map[T]struct{}
	esize    uintptr
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s *Set[T]) Begin() c.Iterator[T] {
	return s.Head()
}

func (s *Set[T]) Head() *it.Iter[T] {
	return it.NewHeadS(s.elements, s.esize)
}

func (s *Set[T]) Revert() *it.Iter[T] {
	return it.NewTailS(s.elements, s.esize)
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

func (s *Set[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Filter(s.Head(), filter))
}

func (s *Set[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Map(s.Head(), by))
}

func (s *Set[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Head(), by)
}

func (s *Set[T]) Len() int {
	return len(s.elements)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *Set[T]) Sort(less func(e1, e2 T) bool) *Set[T] {
	return WrapSet(slice.SortCopy(s.elements, less), s.uniques)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.elements)
}
