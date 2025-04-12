package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
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
	_ collection.Set[int] = (*Set[int])(nil)
	_ collection.Set[int] = Set[int]{}
	_ c.OrderedAll[int]   = Set[int]{}
	_ fmt.Stringer        = (*Set[int])(nil)
	_ fmt.Stringer        = Set[int]{}
)

// All is used to iterate through the collection using `for e := range`.
func (s Set[T]) All(consumer func(T) bool) {
	slice.WalkWhile(s.order, consumer)
}

// IAll is used to iterate through the collection using `for i, e := range`.
func (s Set[T]) IAll(consumer func(int, T) bool) {
	slice.TrackWhile(s.order, consumer)
}

// Loop creates a loop to iterate through the collection.
func (s Set[T]) Loop() loop.Loop[T] {
	return loop.Of(s.order...)
}

// Head creates an iterator to iterate through the collection.
// Deprecated: Head is deprecated. Will be replaced by rance-over function iterator.
func (s Set[T]) Head() slice.Iter[T] {
	return slice.NewHead(s.order)
}

// Tail creates an iterator pointing to the end of the collection
// Deprecated: Tail is deprecated. Will be replaced by rance-over function iterator.
func (s Set[T]) Tail() slice.Iter[T] {
	return slice.NewTail(s.order)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
// Deprecated: First is deprecated. Will be replaced by rance-over function iterator.
func (s Set[T]) First() (*slice.Iter[T], T, bool) {
	iterator := slice.NewHead(s.order)
	return iterator.Crank()
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (s Set[T]) Last() (*slice.Iter[T], T, bool) {
	iterator := slice.NewTail(s.order)
	return iterator.CrankPrev()
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
	return collection.IsEmpty(s)
}

// For applies the 'consumer' function for every element until the consumer returns the c.Break to stop.
func (s Set[T]) For(consumer func(T) error) error {
	return slice.For(s.order, consumer)
}

// ForEach applies the 'consumer' function for every element
func (s Set[T]) ForEach(consumer func(T)) {
	slice.ForEach(s.order, consumer)
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (s Set[T]) Filter(predicate func(T) bool) loop.Loop[T] {
	return loop.Filter(s.Loop(), predicate)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (s Set[T]) Filt(predicate func(T) (bool, error)) breakLoop.Loop[T] {
	return loop.Filt(s.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (s Set[T]) Convert(converter func(T) T) loop.Loop[T] {
	return loop.Convert(s.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (s Set[T]) Conv(converter func(T) (T, error)) breakLoop.Loop[T] {
	return loop.Conv(s.Loop(), converter)
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
