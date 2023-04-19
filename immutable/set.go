package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) *Set[T] {
	internal := map[T]struct{}{}
	for _, e := range elements {
		internal[e] = struct{}{}
	}
	return WrapSet(internal)
}

// WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) *Set[T] {
	return &Set[T]{elements: elements}
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) *Set[T] {
	internal := map[T]struct{}{}
	if elements != nil {
		for e, ok := elements.Next(); ok; e, ok = elements.Next() {
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
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s *Set[T]) Begin() c.Iterator[T] {
	h := s.Head()
	return &h
}

func (s *Set[T]) Head() iter.Key[T, struct{}] {
	var elements map[T]struct{}
	if s != nil {
		elements = s.elements
	}
	return *iter.NewKey(elements)
}

func (s *Set[T]) First() (iter.Key[T, struct{}], T, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (s *Set[T]) Slice() (out []T) {
	if s != nil {
		elements := s.elements
		out = make([]T, 0, len(elements))
		for e := range elements {
			out = append(out, e)
		}
	}
	return out
}

func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) For(walker func(T) error) error {
	if s == nil {
		return nil
	}
	return map_.ForKeys(s.elements, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	if s != nil {
		map_.ForEachKey(s.elements, walker)
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
	return loop.Reduce((&h).Next, by)
}

func (s *Set[T]) Contains(val T) bool {
	if s == nil {
		return false
	}
	_, ok := s.elements[val]
	return ok
}

// Sort transforms to the ordered Set.
func (s *Set[T]) Sort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.Slice, less)
}

func (s *Set[T]) StableSort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s *Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *ordered.Set[T] {
	var elements map[T]struct{}
	if s != nil {
		elements = s.elements
	}
	return ordered.WrapSet(slice.Sort(slice.Clone(s.Slice()), sorter, less), elements)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.Slice())
}
