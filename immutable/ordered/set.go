package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) Set[T] {
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
func WrapSet[T comparable](order []T, elements map[T]struct{}) Set[T] {
	return Set[T]{order: order, elements: elements, esize: notsafe.GetTypeSize[T]()}
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) Set[T] {
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
	_ c.Set[int]   = Set[int]{}
	_ fmt.Stringer = (*Set[int])(nil)
	_ fmt.Stringer = Set[int]{}
)

func (s Set[T]) Begin() c.Iterator[T] {
	return ptr.Of(s.Head())
}

func (s Set[T]) Head() it.ArrayIter[T] {
	return it.NewHeadS(s.order, s.esize)
}

func (s Set[T]) Revert() it.ArrayIter[T] {
	return it.NewTailS(s.order, s.esize)
}

func (s Set[T]) First() (it.ArrayIter[T], T, bool) {
	var (
		iter      = it.NewHeadS(s.order, s.esize)
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (s Set[T]) Last() (it.ArrayIter[T], T, bool) {
	var (
		iter      = it.NewTailS(s.order, s.esize)
		first, ok = iter.Prev()
	)
	return iter, first, ok
}

func (s Set[T]) Collect() []T {
	return slice.Clone(s.order)
}

func (s Set[T]) For(walker func(T) error) error {
	return slice.For(s.order, walker)
}

func (s Set[T]) ForEach(walker func(T)) {
	slice.ForEach(s.order, walker)
}

func (s Set[T]) Filter(filter func(T) bool) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Filter(ptr.Of(s.Head()), filter))
}

func (s Set[T]) Convert(by func(T) T) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Convert(ptr.Of(s.Head()), by))
}

func (s Set[T]) Reduce(by func(T, T) T) T {
	return it.Reduce(ptr.Of(s.Head()), by)
}

func (s Set[T]) Len() int {
	return len(s.order)
}

func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Set[T]) Contains(v T) bool {
	_, ok := s.elements[v]
	return ok
}

func (s Set[T]) Sort(less slice.Less[T]) Set[T] {
	return s.sortBy(sort.Slice, less)
}

func (s Set[T]) StableSort(less slice.Less[T]) Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) Set[T] {
	c := slice.Clone(s.order)
	slice.Sort(c, sorter, less)
	return WrapSet(c, s.elements)
}

func (s Set[T]) String() string {
	return slice.ToString(s.order)
}

func add[T comparable](e T, uniques map[T]struct{}, order []T) []T {
	if _, ok := uniques[e]; !ok {
		order = append(order, e)
		uniques[e] = struct{}{}
	}
	return order
}
