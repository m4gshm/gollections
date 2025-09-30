package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/seq"
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
	_ c.OrderedRange[int] = Set[int]{}
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

// Head returns the first element.
func (s Set[T]) Head() (T, bool) {
	return collection.Head(s)
}

// Tail returns the latest element.
func (s Set[T]) Tail() (T, bool) {
	return slice.Tail(s.order)
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

// ForEach applies the 'consumer' function for every element
func (s Set[T]) ForEach(consumer func(T)) {
	slice.ForEach(s.order, consumer)
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'filter' function
func (s Set[T]) Filter(filter func(T) bool) seq.Seq[T] {
	return collection.Filter(s, filter)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'filter' function
func (s Set[T]) Filt(filter func(T) (bool, error)) seq.SeqE[T] {
	return collection.Filt(s, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (s Set[T]) Convert(converter func(T) T) seq.Seq[T] {
	return collection.Convert(s, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (s Set[T]) Conv(converter func(T) (T, error)) seq.SeqE[T] {
	return collection.Conv(s, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (s Set[T]) Reduce(merge func(T, T) T) T {
	return slice.Reduce(s.order, merge)
}

// HasAny checks whether the set contains an element that satisfies the condition.
func (s Set[T]) HasAny(condition func(T) bool) bool {
	return slice.HasAny(s.order, condition)
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
