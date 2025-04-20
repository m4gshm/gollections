package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) Set[T] {
	return Set[T]{elements: elements}
}

// Set is a collection implementation that provides storage for unique elements, prevents duplication. The elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ collection.Set[int] = (*Set[int])(nil)
	_ collection.Set[int] = Set[int]{}
	_ fmt.Stringer        = (*Set[int])(nil)
	_ fmt.Stringer        = Set[int]{}
)

// All is used to iterate through the collection using `for e := range`.
func (s Set[T]) All(consumer func(T) bool) {
	map_.TrackKeysWhile(s.elements, consumer)
}

// Loop creates a loop to iterate through the collection.
// Deprecated: replaced by the All.
func (s Set[T]) Loop() loop.Loop[T] {
	h := s.Head()
	return (&h).Next
}

// Head creates an iterator to iterate through the collection.
// Deprecated: replaced by the All.
func (s Set[T]) Head() map_.KeyIter[T, struct{}] {
	return map_.NewKeyIter(s.elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
// Deprecated: replaced by the All.
func (s Set[T]) First() (map_.KeyIter[T, struct{}], T, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (s Set[T]) Slice() []T {
	return map_.Keys(s.elements)
}

// Append collects the values to the specified 'out' slice
func (s Set[T]) Append(out []T) []T {
	return map_.AppendKeys(s.elements, out)
}

// Len returns amount of elements
func (s Set[T]) Len() int {
	return len(s.elements)
}

// IsEmpty returns true if the collection is empty
func (s Set[T]) IsEmpty() bool {
	return collection.IsEmpty(s)
}

// For applies the 'consumer' function for the elements until the consumer returns the c.Break to stop.
func (s Set[T]) For(consumer func(T) error) error {

	return map_.ForKeys(s.elements, consumer)
}

// ForEach applies the 'consumer' function for every element
func (s Set[T]) ForEach(consumer func(T)) {
	map_.ForEachKey(s.elements, consumer)
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
	t, _ := map_.Reduce(s.elements, func(t1, t2 T, _, _ struct{}) (t T, out struct{}) {
		return merge(t1, t2), out
	})
	return t
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (s Set[T]) HasAny(predicate func(T) bool) bool {
	return map_.HasAny(s.elements, func(t T, _ struct{}) bool {
		return predicate(t)
	})
}

// Contains checks is the collection contains an element
func (s Set[T]) Contains(element T) (ok bool) {
	if elements := s.elements; elements != nil {
		_, ok = elements[element]
	}
	return ok
}

// Sort transforms to the ordered set with sorted elements
func (s Set[T]) Sort(comparer slice.Comparer[T]) ordered.Set[T] {
	return s.sortBy(slice.Sort, comparer)
}

// StableSort transforms to the ordered set with sorted elements
func (s Set[T]) StableSort(comparer slice.Comparer[T]) ordered.Set[T] {
	return s.sortBy(slice.StableSort, comparer)
}

func (s Set[T]) sortBy(sorter func([]T, slice.Comparer[T]) []T, comparer slice.Comparer[T]) ordered.Set[T] {
	return ordered.WrapSet(sorter(s.Slice(), comparer), s.elements)
}

func (s Set[T]) String() string {
	return slice.ToString(s.Slice())
}
