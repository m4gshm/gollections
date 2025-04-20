package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
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
	_ c.Addable[int]                  = (*Set[int])(nil)
	_ c.AddableNew[int]               = (*Set[int])(nil)
	_ c.AddableAll[c.ForEach[int]]    = (*Set[int])(nil)
	_ c.AddableAllNew[c.ForEach[int]] = (*Set[int])(nil)
	_ c.Deleteable[int]               = (*Set[int])(nil)
	_ c.DeleteableVerify[int]         = (*Set[int])(nil)
	_ c.OrderedRange[int]             = (*Set[int])(nil)
	_ collection.Set[int]             = (*Set[int])(nil)
	_ fmt.Stringer                    = (*Set[int])(nil)
)

// All is used to iterate through the collection using `for e := range`.
func (s *Set[T]) All(consumer func(T) bool) {
	if s != nil {
		slice.WalkWhile(*s.order, consumer)
	}
}

// IAll is used to iterate through the collection using `for i, e := range`.
func (s *Set[T]) IAll(consumer func(int, T) bool) {
	if s != nil {
		slice.TrackWhile(*s.order, consumer)
	}
}

// Loop creates a loop to iterate through the collection.
// Deprecated: replaced by the All.
func (s *Set[T]) Loop() loop.Loop[T] {
	if s == nil {
		return nil
	}
	return loop.Of((*s.order)...)
}

// IterEdit creates iterator that can delete iterable elements
func (s *Set[T]) IterEdit() c.DelIterator[T] {
	h := s.Head()
	return &h
}

// Head creates an iterator to iterate through the collection.
// Deprecated: replaced by the All.
func (s *Set[T]) Head() SetIter[T] {
	var elements *[]T
	if s != nil {
		elements = s.order
	}
	return NewSetIter(elements, s.DeleteOne)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
// Deprecated: replaced by the All.
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
	return collection.IsEmpty(s)
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
func (s *Set[T]) AddAll(other c.ForEach[T]) {
	if !(s == nil || other == nil) {
		other.ForEach(s.AddOne)
	}
}

// AddAllNew inserts elements from the "other" collection if they are not contained in the collection
func (s *Set[T]) AddAllNew(other c.ForEach[T]) (ok bool) {
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

// For applies the 'consumer' function for the elements until the consumer returns the c.Break to stop.
func (s *Set[T]) For(consumer func(T) error) error {
	if s == nil {
		return nil
	}
	order := s.order
	if order == nil {
		return nil
	}
	return slice.For(*order, consumer)
}

// ForEach applies the 'consumer' function for every element
func (s *Set[T]) ForEach(consumer func(T)) {
	if s != nil {
		if order := s.order; order != nil {
			slice.ForEach(*order, consumer)
		}
	}
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filter(predicate func(T) bool) loop.Loop[T] {
	h := s.Head()
	return loop.Filter(h.Next, predicate)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (s *Set[T]) Filt(predicate func(T) (bool, error)) breakLoop.Loop[T] {
	return loop.Filt(s.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (s *Set[T]) Convert(converter func(T) T) loop.Loop[T] {
	return loop.Convert(s.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (s *Set[T]) Conv(converter func(T) (T, error)) breakLoop.Loop[T] {
	return loop.Conv(s.Loop(), converter)
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
