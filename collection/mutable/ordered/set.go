package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	iter "github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
)

// WrapSet creates a set using a map and an order slice as the internal storage.
func WrapSet[T comparable](elements []T, uniques map[T]int) *Set[T] {
	return &Set[T]{order: &elements, elements: uniques}
}

// Set is a collection implementation that provides element uniqueness and access order. Elements must be comparable.
type Set[T comparable] struct {
	order    *[]T
	elements map[T]int
}

var (
	_ c.Addable[int]                      = (*Set[int])(nil)
	_ c.AddableNew[int]                   = (*Set[int])(nil)
	_ c.AddableAll[c.ForEachLoop[int]]    = (*Set[int])(nil)
	_ c.AddableAllNew[c.ForEachLoop[int]] = (*Set[int])(nil)
	_ c.Deleteable[int]                   = (*Set[int])(nil)
	_ c.DeleteableVerify[int]             = (*Set[int])(nil)
	_ collection.Set[int, *SetIter[int]]  = (*Set[int])(nil)
	_ fmt.Stringer                        = (*Set[int])(nil)
)

// Iter creates an iterator and returns as interface
func (s *Set[T]) Iter() *SetIter[T] {
	h := s.Head()
	return &h
}

// IterEdit creates iterator that can delete iterable elements
func (s *Set[T]) IterEdit() c.DelIterator[T] {
	h := s.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (s *Set[T]) Head() SetIter[T] {
	var elements *[]T
	if s != nil {
		elements = s.order
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
		if order := s.order; order != nil {
			out = slice.Clone(*s.order)
		}
	}
	return out
}

// Append collects the values to the specified 'out' slice
func (s *Set[T]) Append(out []T) []T {
	if s != nil {
		if order := s.order; order != nil {
			out = append(out, (*s.order)...)
		}
	}
	return out
}

// Clone returns copy of the collection
func (s *Set[T]) Clone() *Set[T] {
	var (
		elements []T
		uniques  map[T]int
	)
	if s != nil {
		if order := s.order; order != nil {
			elements = slice.Clone(*s.order)
		}
		uniques = map_.Clone(s.elements)
	}
	return WrapSet(elements, uniques)
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
	if order := s.order; order != nil {
		return len(*s.order)
	}
	return 0
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
	s.AddNew(elements...)
}

// AddOne adds an element in the collection
func (s *Set[T]) AddOne(element T) {
	s.AddOneNew(element)
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
			elements = map[T]int{}
			s.elements = elements
			s.order = &[]T{}
		}
		if ok = !s.Contains(element); ok {
			order := s.order
			if order != nil {
				elements[element] = len(*order)
				*(s.order) = append(*order, element)
			}
		}
	}
	return ok
}

// AddAll inserts all elements from the "other" collection
func (s *Set[T]) AddAll(other c.ForEachLoop[T]) {
	if !(s == nil || other == nil) {
		other.ForEach(s.AddOne)
	}
}

// AddAllNew inserts elements from the "other" collection if they are not contained in the collection
func (s *Set[T]) AddAllNew(other c.ForEachLoop[T]) (ok bool) {
	if !(s == nil || other == nil) {
		other.ForEach(func(v T) { ok = s.AddOneNew(v) || ok })
	}
	return ok
}

// Delete removes elements from the collection
func (s *Set[T]) Delete(elements ...T) {
	s.DeleteActual(elements...)
}

// DeleteOne removes an element from the collection
func (s *Set[T]) DeleteOne(v T) {
	s.DeleteActualOne(v)
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
func (s *Set[T]) DeleteActualOne(element T) bool {
	if !(s == nil || s.elements == nil) {
		elements := s.elements
		if pos, ok := elements[element]; ok {
			delete(elements, element)
			//todo: need optimize
			order := s.order
			ne := slice.Delete(*order, pos)
			for i := pos; i < len(ne); i++ {
				elements[ne[i]]--
			}
			*(s.order) = ne
			return true
		}
	}
	return false
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (s *Set[T]) For(walker func(T) error) error {
	if s == nil {
		return nil
	}
	order := s.order
	if order == nil {
		return nil
	}
	return slice.For(*order, walker)
}

// ForEach applies the 'walker' function for every element
func (s *Set[T]) ForEach(walker func(T)) {
	if s != nil {
		if order := s.order; order != nil {
			slice.ForEach(*order, walker)
		}
	}
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filter(predicate func(T) bool) stream.Iter[T] {
	h := s.Head()
	return stream.New(iter.Filter(h.Next, predicate).Next)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filt(predicate func(T) (bool, error)) breakStream.Iter[T] {
	h := s.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (s *Set[T]) Convert(converter func(T) T) stream.Iter[T] {
	return collection.Convert(s, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (s *Set[T]) Conv(converter func(T) (T, error)) breakStream.Iter[T] {
	return collection.Conv(s, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s *Set[T]) Reduce(merge func(T, T) T) (t T) {
	if s != nil {
		if order := s.order; order != nil {
			t = slice.Reduce(*order, merge)
		}
	}
	return t
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (s *Set[K]) HasAny(predicate func(K) bool) bool {
	if s != nil {
		if order := s.order; order != nil {
			return slice.HasAny(*order, predicate)
		}
	}
	return false
}

// Sort sorts the elements
func (s *Set[T]) Sort(comparer slice.Comparer[T]) *Set[T] {
	return s.sortBy(slice.Sort, comparer)
}

// StableSort sorts the elements
func (s *Set[T]) StableSort(comparer slice.Comparer[T]) *Set[T] {
	return s.sortBy(slice.StableSort, comparer)
}

func (s *Set[T]) sortBy(sorter func([]T, slice.Comparer[T]) []T, comparer slice.Comparer[T]) *Set[T] {
	if s != nil {
		if order := s.order; order != nil {
			sorter(*order, comparer)
			for i, v := range *order {
				s.elements[v] = i
			}
		}
	}
	return s
}

func (s *Set[T]) String() string {
	var elements []T
	if s != nil {
		if order := s.order; order != nil {
			elements = *order
		}
	}
	return slice.ToString(elements)
}

func addToSet[T comparable](e T, uniques map[T]int, order []T, pos int) ([]T, int) {
	if _, ok := uniques[e]; !ok {
		order = append(order, e)
		uniques[e] = pos
		pos++
	}
	return order, pos
}
