package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// WrapVector instantiates Vector using a slise as internal storage.
func WrapVector[T any](elements []T) Vector[T] {
	return Vector[T]{elements: elements}
}

// Vector is a collection implementation that provides elements order and index access.
type Vector[T any] struct {
	elements []T
}

var (
	_ collection.Vector[any] = (*Vector[any])(nil)
	_ collection.Vector[any] = Vector[any]{}
	_ fmt.Stringer           = (*Vector[any])(nil)
	_ fmt.Stringer           = Vector[any]{}
)

// Loop creates a loop to iterating through elements.
func (v Vector[T]) Loop() loop.Loop[T] {
	return v.Head().Next
}

// Head creates an iterator and returns as implementation type value
func (v Vector[T]) Head() *slice.Iter[T] {
	return slice.NewHead(v.elements)
}

// Tail creates an iterator pointing to the end of the collection
func (v Vector[T]) Tail() *slice.Iter[T] {
	return slice.NewTail(v.elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (v Vector[T]) First() (*slice.Iter[T], T, bool) {
	return slice.NewHead(v.elements).Crank()
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (v Vector[T]) Last() (*slice.Iter[T], T, bool) {
	return slice.NewTail(v.elements).CrankPrev()
}

// Slice collects the elements to a slice
func (v Vector[T]) Slice() []T {
	if elements := v.elements; elements != nil {
		return slice.Clone(elements)
	}
	return nil
}

// Append collects the values to the specified 'out' slice
func (v Vector[T]) Append(out []T) []T {
	if elements := v.elements; elements != nil {
		return append(out, elements...)
	}
	return out
}

// Len returns amount of elements
func (v Vector[T]) Len() int {
	return notsafe.GetLen(v.elements)
}

// IsEmpty returns true if the collection is empty
func (v Vector[T]) IsEmpty() bool {
	return v.Len() == 0
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the vector len, returns zero T and false in the second result
func (v Vector[T]) Get(index int) (out T, ok bool) {
	return slice.Gett(v.elements, index)
}

// Track applies the 'tracker' function for elements. Return the c.Break to stop.
func (v Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(v.elements, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (v Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(v.elements, tracker)

}

// For applies the 'walker' function for the elements. Return the c.Break to stop.
func (v Vector[T]) For(walker func(T) error) error {
	return slice.For(v.elements, walker)
}

// ForEach applies the 'walker' function for every element
func (v Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(v.elements, walker)
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (v Vector[T]) Filter(filter func(T) bool) loop.Loop[T] {
	return loop.Filter(v.Loop(), filter)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (v Vector[T]) Filt(predicate func(T) (bool, error)) breakLoop.Loop[T] {
	return loop.Filt(v.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (v Vector[T]) Convert(converter func(T) T) loop.Loop[T] {
	return loop.Convert(v.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (v Vector[T]) Conv(converter func(T) (T, error)) breakLoop.Loop[T] {
	return loop.Conv(v.Loop(), converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (v Vector[T]) Reduce(merge func(T, T) T) T {
	return slice.Reduce(v.elements, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (v Vector[T]) HasAny(predicate func(T) bool) bool {
	return slice.HasAny(v.elements, predicate)
}

// Sort returns a sorted clone of the Vector
func (v Vector[T]) Sort(comparer slice.Comparer[T]) Vector[T] {
	return v.sortBy(slice.Sort, comparer)
}

// StableSort returns a stable sorted clone of the Vector
func (v Vector[T]) StableSort(comparer slice.Comparer[T]) Vector[T] {
	return v.sortBy(slice.StableSort, comparer)
}

func (v Vector[T]) sortBy(sorter func([]T, slice.Comparer[T]) []T, comparer slice.Comparer[T]) Vector[T] {
	return WrapVector(sorter(slice.Clone(v.elements), comparer))
}

func (v Vector[T]) String() string {
	return slice.ToString(v.elements)
}
