package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) *Set[T] {
	var (
		uniques = map[T]struct{}{}
		order   []T
	)
	for _, e := range elements {
		order = add(e, uniques, order)
	}
	return WrapSet(order, uniques)
}

// WrapSet creates a set using a map and an order slice as the internal storage.
func WrapSet[T comparable](order []T, elements map[T]struct{}) *Set[T] {
	return &Set[T]{order: order, elements: elements, esize: notsafe.GetTypeSize[T]()}
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) *Set[T] {
	var (
		uniques = map[T]struct{}{}
		order   []T
	)
	for {
		if e, ok := elements.Next(); !ok {
			break
		} else {
			order = add(e, uniques, order)
		}
	}
	return WrapSet(order, uniques)
}

// Set is the Collection implementation that provides element uniqueness and access order. The elements must be comparable.
type Set[T comparable] struct {
	order    []T
	elements map[T]struct{}
	esize    uintptr
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s Set[T]) P() *Set[T] {
	return &s
}

func (s *Set[T]) Begin() c.Iterator[T] {
	h := s.Head()
	return &h
}

func (s *Set[T]) Head() iter.ArrayIter[T] {
	var (
		order []T
		esize uintptr
	)
	if s != nil {
		order = s.order
		esize = s.esize
	}
	return iter.NewHeadS(order, esize)
}

func (s *Set[T]) Revert() iter.ArrayIter[T] {
	var (
		order []T
		esize uintptr
	)
	if s != nil {
		order = s.order
		esize = s.esize
	}
	return iter.NewTailS(order, esize)
}

func (s *Set[T]) First() (iter.ArrayIter[T], T, bool) {
	var (
		order []T
		esize uintptr
	)
	if s != nil {
		order = s.order
		esize = s.esize
	}
	var (
		iterator  = iter.NewHeadS(order, esize)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (s *Set[T]) Last() (iter.ArrayIter[T], T, bool) {
	var (
		order []T
		esize uintptr
	)
	if s != nil {
		order = s.order
		esize = s.esize
	}
	var (
		iterator  = iter.NewTailS(order, esize)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

func (s *Set[T]) Slice() (out []T) {
	if s != nil {
		out = slice.Clone(s.order)
	}
	return out
}

func (s *Set[T]) For(walker func(T) error) error {
	if s == nil {
		return nil
	}
	return slice.For(s.order, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	if s != nil {
		slice.ForEach(s.order, walker)
	}
}

func (s *Set[T]) Filter(filter func(T) bool) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Filter(h, h.Next, filter))
}

func (s *Set[T]) Convert(by func(T) T) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Convert(h, h.Next, by))
}

func (s *Set[T]) Reduce(by func(T, T) T) T {
	h := s.Head()
	return loop.Reduce(h.Next, by)
}

func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.order)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Contains(v T) (ok bool) {
	if s != nil {
		_, ok = s.elements[v]
	}
	return
}

func (s *Set[T]) Sort(less slice.Less[T]) *Set[T] {
	return s.sortBy(sort.Slice, less)
}

func (s *Set[T]) StableSort(less slice.Less[T]) *Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s *Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *Set[T] {
	var order []T
	if s != nil {
		order = slice.Clone(s.order)
		slice.Sort(order, sorter, less)
	}
	return WrapSet(order, s.elements)
}

func (s *Set[T]) String() string {
	if s == nil {
		return ""
	}
	return slice.ToString(s.order)
}

func add[T comparable](e T, uniques map[T]struct{}, order []T) []T {
	if _, ok := uniques[e]; !ok {
		order = append(order, e)
		uniques[e] = struct{}{}
	}
	return order
}
