package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
)

// WrapSet creates a set using a map and an order slice as the internal storage.
func WrapSet[T comparable](order []T, elements map[T]struct{}) Set[T] {
	return Set[T]{order: order, elements: elements}
}

// Set is a collection implementation that provides storage for unique elements, prevents duplication, and guarantees access order. The elements must be comparable.
type Set[T comparable] struct {
	order    []T
	elements map[T]struct{}
}

var (
	_ collection.Set[int, *slice.Iter[int]] = (*Set[int])(nil)
	_ collection.Set[int, *slice.Iter[int]] = Set[int]{}
	_ fmt.Stringer                         = (*Set[int])(nil)
	_ fmt.Stringer                         = Set[int]{}
)

// Iter creates an iterator and returns as interface
func (s Set[T]) Iter() *slice.Iter[T] {
	h := s.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (s Set[T]) Head() slice.Iter[T] {
	return slice.NewHead(s.order)
}

// Tail creates an iterator pointing to the end of the collection
func (s Set[T]) Tail() slice.Iter[T] {
	return slice.NewTail(s.order)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (s Set[T]) First() (slice.Iter[T], T, bool) {
	var (
		iterator  = slice.NewHead(s.order)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (s Set[T]) Last() (slice.Iter[T], T, bool) {
	var (
		iterator  = slice.NewTail(s.order)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (s Set[T]) Slice() []T {
	return slice.Clone(s.order)
}

// Append collects the values to the specified 'out' slice
func (s Set[T]) Append(out []T) []T {
	return append(out, s.order...)
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

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (s Set[T]) Filter(predicate func(T) bool) stream.Iter[T] {
	h := s.Head()
	return stream.New(loop.Filter(h.Next, predicate).Next)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (s Set[T]) Filt(predicate func(T) (bool, error)) breakStream.Iter[T] {
	h := s.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (s Set[T]) Convert(converter func(T) T) stream.Iter[T] {
	return collection.Convert(s, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (s Set[T]) Conv(converter func(T) (T, error)) breakStream.Iter[T] {
	return collection.Conv(s, converter)
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
func (s Set[T]) Sort(comparer slice.Comparer[T]) Set[T] {
	return s.sortBy(slice.Sort, comparer)
}

// StableSort sorts the elements
func (s Set[T]) StableSort(comparer slice.Comparer[T]) Set[T] {
	return s.sortBy(slice.StableSort, comparer)
}

func (s Set[T]) sortBy(sorter func([]T, slice.Comparer[T]) []T, comparer slice.Comparer[T]) Set[T] {
	return WrapSet(sorter(slice.Clone(s.order), comparer), s.elements)
}

func (s Set[T]) String() string {
	return slice.ToString(s.order)
}

func addToSet[T comparable](e T, uniques map[T]struct{}, order []T) []T {
	if _, ok := uniques[e]; !ok {
		order = append(order, e)
		uniques[e] = struct{}{}
	}
	return order
}
