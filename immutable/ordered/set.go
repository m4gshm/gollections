package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/notsafe"
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
	for e, ok := elements.Next(); ok; e, ok = elements.Next() {
		order = add(e, uniques, order)
	}
	return WrapSet(order, uniques)
}

// Set is a collection that provides storage for unique elements, prevents duplication, and guarantees access order. The elements must be comparable.
type Set[T comparable] struct {
	order    []T
	elements map[T]struct{}
	esize    uintptr
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
	_ c.Set[int]   = Set[int]{}
	_ fmt.Stringer = Set[int]{}
)

// Begin creates iterator
func (s Set[T]) Begin() c.Iterator[T] {
	h := s.Head()
	return &h
}

// Head creates iterator
func (s Set[T]) Head() iter.ArrayIter[T] {
	return iter.NewHeadS(s.order, s.esize)
}

// Tail creates an iterator pointing to the end of the collection
func (s Set[T]) Tail() iter.ArrayIter[T] {
	return iter.NewTailS(s.order, s.esize)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (s Set[T]) First() (iter.ArrayIter[T], T, bool) {
	var (
		iterator  = iter.NewHeadS(s.order, s.esize)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (s Set[T]) Last() (iter.ArrayIter[T], T, bool) {
	var (
		iterator  = iter.NewTailS(s.order, s.esize)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (s Set[T]) Slice() []T {
	return slice.Clone(s.order)
}

// Len returns amount of elements
func (s Set[T]) Len() int {
	return len(s.order)
}

// IsEmpty returns true if the collection is empty
func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

// For applies the 'walker' function for every element. Return the c.ErrBreak to stop.
func (s Set[T]) For(walker func(T) error) error {
	return slice.For(s.order, walker)
}

// ForEach applies the 'walker' function for every element
func (s Set[T]) ForEach(walker func(T)) {
	slice.ForEach(s.order, walker)
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (s Set[T]) Filter(predicate func(T) bool) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Filter(h.Next, predicate).Next)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (s Set[T]) Convert(by func(T) T) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Convert(h.Next, by).Next)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s Set[T]) Reduce(merge func(T, T) T) T {
	return slice.Reduce(s.order, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (s Set[T]) HasAny(predicate func(T) bool) bool {
	return slice.HasAny(s.order, predicate)
}

// Contains checks is the collection contains an element
func (s Set[T]) Contains(element T) (ok bool) {
	if elements := s.elements; elements != nil {
		_, ok = s.elements[element]
	}
	return ok
}

// Sort sorts the elements
func (s Set[T]) Sort(less slice.Less[T]) Set[T] {
	return s.sortBy(sort.Slice, less)
}

// StableSort sorts the elements
func (s Set[T]) StableSort(less slice.Less[T]) Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) Set[T] {
	order := slice.Clone(s.order)
	slice.Sort(order, sorter, less)
	return WrapSet(order, s.elements)
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
