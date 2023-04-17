package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) Set[T] {
	internal := map[T]struct{}{}
	for _, e := range elements {
		internal[e] = struct{}{}
	}
	return WrapSet(internal)
}

// WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) Set[T] {
	return Set[T]{elements: elements}
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) Set[T] {
	internal := map[T]struct{}{}
	for {
		if e, ok := elements.Next(); !ok {
			break
		} else {
			internal[e] = struct{}{}
		}
	}
	return WrapSet(internal)
}

// Set is the Collection implementation that provides the uniqueness of its elements. Elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ c.Set[int]   = Set[int]{}
	_ fmt.Stringer = (*Set[int])(nil)
	_ fmt.Stringer = Set[int]{}
)

func (s Set[T]) Begin() c.Iterator[T] {
	return s.Head()
}

func (s Set[T]) Head() it.Key[T, struct{}] {
	return it.NewKey(s.elements)
}

func (s Set[T]) First() (it.Key[T, struct{}], T, bool) {
	var (
		iter      = s.Head()
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (s Set[T]) Collect() []T {
	elements := s.elements
	out := make([]T, 0, len(elements))
	for e := range elements {
		out = append(out, e)
	}
	return out
}

func (s Set[T]) Slice() []T {
	return s.Collect()
}

func (s Set[T]) Len() int {
	return len(s.elements)
}

func (s Set[T]) IsEmpty() bool {
	return s.IsEmpty()
}

func (s Set[T]) For(walker func(T) error) error {
	return map_.ForKeys(s.elements, walker)
}

func (s Set[T]) ForEach(walker func(T)) {
	map_.ForEachKey(s.elements, walker)
}

func (s Set[T]) Filter(filter func(T) bool) c.Pipe[T, []T] {
	h := s.Head()
	return it.NewPipe[T](it.Filter(h, h.Next, filter))
}

func (s Set[T]) Convert(by func(T) T) c.Pipe[T, []T] {
	h := s.Head()
	return it.NewPipe[T](it.Convert(h, h.Next, by))
}

func (s Set[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(s.Head().Next, by)
}

func (s Set[T]) Contains(val T) bool {
	_, ok := s.elements[val]
	return ok
}

// Sort transforms to the ordered Set.
func (s Set[T]) Sort(less slice.Less[T]) ordered.Set[T] {
	return s.sortBy(sort.Slice, less)
}

func (s Set[T]) StableSort(less slice.Less[T]) ordered.Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) ordered.Set[T] {
	c := slice.Clone(s.Slice())
	slice.Sort(c, sorter, less)
	return ordered.WrapSet(c, s.elements)
}

func (s Set[T]) String() string {
	return slice.ToString(s.Slice())
}
