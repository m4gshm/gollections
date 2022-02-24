package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

//NewSet creates the Set and copies elements to it.
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

//WrapSet creates a set using a map and an order slice as the internal storage.
func WrapSet[T comparable](order []T, elements map[T]struct{}) *Set[T] {
	return &Set[T]{order: order, elements: elements, esize: notsafe.GetTypeSize[T]()}
}

//Set is the Collection implementation that provides element uniqueness and access order. The elements must be comparable.
type Set[T comparable] struct {
	order    []T
	elements map[T]struct{}
	esize    uintptr
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s *Set[T]) Begin() c.Iterator[T] {
	iter := s.Head()
	return &iter
}

func (s *Set[T]) Head() it.Iter[T] {
	return it.NewHeadS(s.order, s.esize)
}

func (s *Set[T]) Revert() it.Iter[T] {
	return it.NewTailS(s.order, s.esize)
}

func (s *Set[T]) First() (it.Iter[T], T, bool) {
	var (
		iter      = it.NewHeadS(s.order, s.esize)
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (s *Set[T]) Last() (it.Iter[T], T, bool) {
	var (
		iter      = it.NewTailS(s.order, s.esize)
		first, ok = iter.Prev()
	)
	return iter, first, ok
}

func (s *Set[T]) Collect() []T {
	return slice.Copy(s.order)
}

func (s *Set[T]) For(walker func(T) error) error {
	return slice.For(s.order, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	slice.ForEach(s.order, walker)
}

func (s *Set[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	iter := s.Head()
	return it.NewPipe[T](it.Filter(&iter, filter))
}

func (s *Set[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	iter := s.Head()
	return it.NewPipe[T](it.Map(&iter, by))
}

func (s *Set[T]) Reduce(by op.Binary[T]) T {
	iter := s.Head()
	return it.Reduce(&iter, by)
}

func (s *Set[T]) Len() int {
	return len(s.order)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.elements[v]
	return ok
}

func (s *Set[T]) Sort(less func(e1, e2 T) bool) *Set[T] {
	return WrapSet(slice.SortCopy(s.order, less), s.elements)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.order)
}
