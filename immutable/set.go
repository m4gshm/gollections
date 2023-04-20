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

// Set is a collection that provides storage for unique elements, prevents duplication. The elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

// Begin creates iterator
func (s *Set[T]) Begin() c.Iterator[T] {
	h := s.Head()
	return &h
}

// Head creates iterator
func (s *Set[T]) Head() iter.Key[T, struct{}] {
	var elements map[T]struct{}
	if s != nil {
		elements = s.elements
	}
	return *iter.NewKey(elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (s *Set[T]) First() (iter.Key[T, struct{}], T, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (s *Set[T]) Slice() (out []T) {
	if s != nil {
		out = map_.Keys(s.elements)
	}
	return out
}

// Len returns amount of elements
func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

// IsEmpty returns true if the collection is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (s *Set[T]) For(walker func(T) error) error {
	if s == nil {
		return nil
	}
	return map_.ForKeys(s.elements, walker)
}

// ForEach applies the 'walker' function for every element
func (s *Set[T]) ForEach(walker func(T)) {
	if s != nil {
		map_.ForEachKey(s.elements, walker)
	}
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filter(filter func(T) bool) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Filter(h, h.Next, filter))
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (s *Set[T]) Convert(converter func(T) T) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Convert(h, h.Next, converter))
}

// Reduce reduces the elements into an one using the 'merge' function
func (s *Set[T]) Reduce(by func(T, T) T) T {
	h := s.Head()
	return loop.Reduce((&h).Next, by)
}

// Contains checks is the collection contains an element
func (s *Set[T]) Contains(element T) bool {
	if s == nil {
		return false
	}
	_, ok := s.elements[element]
	return ok
}

// Sort transforms to the ordered set with sorted elements
func (s *Set[T]) Sort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.Slice, less)
}

// StableSort transforms to the ordered set with sorted elements
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
