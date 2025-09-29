package mutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

// WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) *Set[T] {
	return &Set[T]{elements: elements}
}

// Set is a collection implementation that provides element uniqueness. The elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ c.Addable[int]                = (*Set[int])(nil)
	_ c.AddableNew[int]             = (*Set[int])(nil)
	_ c.AddableAll[seq.Seq[int]]    = (*Set[int])(nil)
	_ c.AddableAllNew[seq.Seq[int]] = (*Set[int])(nil)
	_ c.Deleteable[int]             = (*Set[int])(nil)
	_ c.DeleteableVerify[int]       = (*Set[int])(nil)
	_ collection.Set[int]           = (*Set[int])(nil)
	_ fmt.Stringer                  = (*Set[int])(nil)
)

// All is used to iterate through the collection using `for e := range`.
func (s *Set[T]) All(consumer func(T) bool) {
	if s == nil {
		return
	}
	for v := range s.elements {
		if !consumer(v) {
			return
		}
	}

}

// Head returns the first element.
func (s *Set[T]) Head() (T, bool) {
	return collection.Head(s)
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

// Append collects the values to the specified 'out' slice
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
	return collection.IsEmpty(s)
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

// AddAll inserts all elements from the "other" sequence.
func (s *Set[T]) AddAll(elements seq.Seq[T]) {
	if !(s == nil || elements == nil) {
		seq.ForEach(elements, s.AddOne)
	}
}

// AddAllNew inserts elements from the "other" sequence if they are not contained in the collection.
func (s *Set[T]) AddAllNew(other seq.Seq[T]) (ok bool) {
	if !(s == nil || other == nil) {
		seq.ForEach(other, func(element T) { ok = s.AddOneNew(element) || ok })
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

// For applies the 'consumer' function for the elements until the consumer returns the c.Break to stop.
func (s *Set[T]) For(consumer func(T) error) error {
	if s == nil {
		return nil
	}
	return map_.ForKeys(s.elements, consumer)
}

// ForEach applies the 'consumer' function for every element
func (s *Set[T]) ForEach(consumer func(T)) {
	if s != nil {
		map_.ForEachKey(s.elements, consumer)
	}
}

// Filter returns a seq that checks elements by the 'filter' function and returns successful ones.
func (s *Set[T]) Filter(predicate func(T) bool) seq.Seq[T] {
	return collection.Filter(s, predicate)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filt(predicate func(T) (bool, error)) seq.SeqE[T] {
	return collection.Filt(s, predicate)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (s *Set[T]) Convert(converter func(T) T) seq.Seq[T] {
	return collection.Convert(s, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (s *Set[T]) Conv(converter func(T) (T, error)) seq.SeqE[T] {
	return collection.Conv(s, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s *Set[T]) Reduce(merge func(T, T) T) (t T) {
	if s != nil {
		t, _ = map_.Reduce(s.elements, func(t1, t2 T, _, _ struct{}) (t T, out struct{}) {
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
func (s *Set[T]) Sort(comparer slice.Comparer[T]) *ordered.Set[T] {
	if s != nil {
		return ordered.NewSet(slice.Sort(s.Slice(), comparer)...)
	}
	return nil
}

// StableSort transforms to the ordered Set contains sorted elements
func (s *Set[T]) StableSort(comparer slice.Comparer[T]) *ordered.Set[T] {
	if s != nil {
		return ordered.NewSet(slice.StableSort(s.Slice(), comparer)...)
	}
	return nil
}

func (s *Set[T]) String() string {
	return slice.ToString(s.Slice())
}
