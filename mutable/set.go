package mutable

import (
	"fmt"
	"sort"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
)

// NewSetCap creates a set with a predefined capacity.
func NewSetCap[T comparable](capacity int) *Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) *Set[T] {
	internal := make(map[T]struct{}, len(elements))
	for _, v := range elements {
		internal[v] = struct{}{}
	}
	return WrapSet(internal)
}

// SetFromLoop creates a set with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func SetFromLoop[T comparable](next func() (T, bool)) *Set[T] {
	internal := map[T]struct{}{}
	for e, ok := next(); ok; e, ok = next() {
		internal[e] = struct{}{}
	}
	return WrapSet(internal)
}

// WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) *Set[T] {
	return &Set[T]{elements: elements}
}

// Set is the Collection implementation that provides element uniqueness. The elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ c.Addable[int]                      = (*Set[int])(nil)
	_ c.AddableNew[int]                   = (*Set[int])(nil)
	_ c.AddableAll[c.ForEachLoop[int]]    = (*Set[int])(nil)
	_ c.AddableAllNew[c.ForEachLoop[int]] = (*Set[int])(nil)
	_ c.Deleteable[int]                   = (*Set[int])(nil)
	_ c.DeleteableVerify[int]             = (*Set[int])(nil)
	_ collection.Set[int]                 = (*Set[int])(nil)
	_ loop.Looper[int, *SetIter[int]]     = (*Set[int])(nil)
	_ fmt.Stringer                        = (*Set[int])(nil)
)

// Iter creates an iterator
func (s *Set[T]) Iter() c.Iterator[T] {
	h := s.Head()
	return &h
}

func (s *Set[T]) Loop() *SetIter[T] {
	h := s.Head()
	return &h
}

// IterEdit creates iterator that can delete iterable elements
func (s *Set[T]) IterEdit() c.DelIterator[T] {
	h := s.Head()
	return &h
}

// Head creates an iterator value object
func (s *Set[T]) Head() SetIter[T] {
	var elements map[T]struct{}
	if s != nil {
		elements = s.elements
	}
	return NewSetIter(elements, s.DeleteOne)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (s *Set[T]) First() (SetIter[T], T, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (s *Set[T]) Slice() (out []T) {
	if s != nil {
		if elements := s.elements; elements != nil {
			out = map_.Keys(elements)
		}
	}
	return out
}

func (s *Set[T]) Append(out []T) []T {
	if s != nil {
		if elements := s.elements; elements != nil {
			out = map_.AppendKeys(s.elements, out)
		}
	}
	return out
}

// Clone returns copy of the collection
func (s *Set[T]) Clone() *Set[T] {
	var elements map[T]struct{}
	if s != nil {
		elements = map_.Clone(s.elements)
	}
	return WrapSet(elements)
}

// IsEmpty returns true if the collection is empty
func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Len returns amount of the elements
func (s *Set[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.elements)
}

// Contains checks if the collection contains an element
func (s *Set[T]) Contains(element T) (ok bool) {
	if s != nil {
		_, ok = s.elements[element]
	}
	return ok
}

// Add adds elements in the collection
func (s *Set[T]) Add(elements ...T) {
	if s != nil {
		for _, element := range elements {
			s.AddOne(element)
		}
	}
}

// AddOne adds an element in the collection
func (s *Set[T]) AddOne(element T) {
	if s == nil {
		return
	} else if elements := s.elements; elements == nil {
		s.elements = map[T]struct{}{element: {}}
	} else {
		elements[element] = struct{}{}
	}
}

// AddNew inserts elements if they are not contained in the collection
func (s *Set[T]) AddNew(elements ...T) bool {
	if s == nil {
		return false
	}
	ok := false
	for _, element := range elements {
		ok = s.AddOneNew(element) || ok
	}
	return ok
}

// AddOneNew inserts an element if it is not contained in the collection
func (s *Set[T]) AddOneNew(element T) (ok bool) {
	if s != nil {
		elements := s.elements
		if elements == nil {
			elements = map[T]struct{}{}
			s.elements = elements
		}
		if ok = !s.Contains(element); ok {
			elements[element] = struct{}{}
		}
	}
	return ok
}

// AddAll inserts all elements from the "other" collection
func (s *Set[T]) AddAll(elements c.ForEachLoop[T]) {
	if !(s == nil || elements == nil) {
		elements.ForEach(s.AddOne)
	}
}

// AddAllNew inserts elements from the "other" collection if they are not contained in the collection
func (s *Set[T]) AddAllNew(other c.ForEachLoop[T]) (ok bool) {
	if !(s == nil || other == nil) {
		other.ForEach(func(element T) { ok = s.AddOneNew(element) || ok })
	}
	return ok
}

// Delete removes elements from the collection
func (s *Set[T]) Delete(elements ...T) {
	if s != nil {
		for _, element := range elements {
			s.DeleteOne(element)
		}
	}
}

// DeleteOne removes an element from the collection
func (s *Set[T]) DeleteOne(element T) {
	if s != nil {
		delete(s.elements, element)
	}
}

// DeleteActual removes elements only if they are contained in the collection
func (s *Set[T]) DeleteActual(elements ...T) bool {
	if s == nil {
		return false
	}
	ok := false
	for i := range elements {
		ok = s.DeleteActualOne(elements[i]) || ok
	}
	return ok
}

// DeleteActualOne removes an element only if it is contained in the collection
func (s *Set[T]) DeleteActualOne(element T) (ok bool) {
	if !(s == nil || s.elements == nil) {
		if _, ok = s.elements[element]; ok {
			delete(s.elements, element)
		}
	}
	return ok
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

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filter(predicate func(T) bool) stream.Iter[T] {
	h := s.Head()
	return stream.New(loop.Filter(h.Next, predicate).Next)
}

// Filt returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (s Set[T]) Filt(predicate func(T) (bool, error)) breakStream.Iter[T] {
	h := s.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (s *Set[T]) Convert(converter func(T) T) stream.Iter[T] {
	return collection.Convert(s, converter)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (s *Set[T]) Conv(converter func(T) (T, error)) breakStream.Iter[T] {
	return collection.Conv(s, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s *Set[T]) Reduce(merge func(T, T) T) (t T) {
	if s != nil {
		t, _ = map_.Reduce(s.elements, func(t1 T, _ struct{}, t2 T, _ struct{}) (t T, out struct{}) {
			return merge(t1, t2), out
		})
	}
	return t
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (s *Set[K]) HasAny(predicate func(K) bool) bool {
	if s != nil {
		return map_.HasAny(s.elements, func(k K, _ struct{}) bool {
			return predicate(k)
		})
	}
	return false
}

// Sort transforms to the ordered Set contains sorted elements
func (s *Set[T]) Sort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.Slice, less)
}

// StableSort transforms to the ordered Set contains sorted elements
func (s *Set[T]) StableSort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s *Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *ordered.Set[T] {
	var sortedElements []T
	if s != nil {
		sortedElements = slice.Clone(s.Slice())
		slice.Sort(sortedElements, sorter, less)
	}
	return ordered.NewSet(sortedElements)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.Slice())
}
