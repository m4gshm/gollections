package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
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

// Head returns the first element.
func (s Set[T]) Head() (T, bool) {
	return collection.Head(s)
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

// ForEach applies the 'consumer' function for every element
func (s Set[T]) ForEach(consumer func(T)) {
	map_.ForEachKey(s.elements, consumer)
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'filter' function
func (s Set[T]) Filter(filter func(T) bool) seq.Seq[T] {
	return collection.Filter(s, filter)
}

// Filt returns an errorable seq consisting of elements that satisfy the condition of the 'filter' function
func (s Set[T]) Filt(filter func(T) (bool, error)) seq.SeqE[T] {
	return collection.Filt(s, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (s Set[T]) Convert(converter func(T) T) seq.Seq[T] {
	return collection.Convert(s, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func (s Set[T]) Conv(converter func(T) (T, error)) seq.SeqE[T] {
	return collection.Conv(s, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s Set[T]) Reduce(merge func(T, T) T) T {
	t, _ := map_.Reduce(s.elements, func(t1, t2 T, _, _ struct{}) (t T, out struct{}) {
		return merge(t1, t2), out
	})
	return t
}

// HasAny checks whether the set contains an element that satisfies the condition.
func (s Set[T]) HasAny(condition func(T) bool) bool {
	return map_.HasAny(s.elements, func(t T, _ struct{}) bool {
		return condition(t)
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
