package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
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
	_ collection.Vector[any]             = (*Vector[any])(nil)
	_ collection.Vector[any]             = Vector[any]{}
	_ loop.Looper[any, *slice.Iter[any]] = (*Vector[any])(nil)
	_ loop.Looper[any, *slice.Iter[any]] = Vector[any]{}
	_ fmt.Stringer                       = (*Vector[any])(nil)
	_ fmt.Stringer                       = Vector[any]{}
)

func (v Vector[T]) All(yield func(T) bool) {
	slice.All(v.elements, yield)
}

// Iter creates an iterator and returns as interface
func (v Vector[T]) Iter() c.Iterator[T] {
	h := v.Head()
	return &h
}

// Loop creates an iterator and returns as implementation type reference
func (v Vector[T]) Loop() *slice.Iter[T] {
	h := v.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (v Vector[T]) Head() slice.Iter[T] {
	return slice.NewHead(v.elements)
}

// Tail creates an iterator pointing to the end of the collection
func (v Vector[T]) Tail() slice.Iter[T] {
	return slice.NewTail(v.elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (v Vector[T]) First() (slice.Iter[T], T, bool) {
	var (
		iterator  = slice.NewHead(v.elements)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (v Vector[T]) Last() (slice.Iter[T], T, bool) {
	var (
		iterator  = slice.NewTail(v.elements)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
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

// Track applies the 'tracker' function for elements. Return the c.ErrBreak to stop.
func (v Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(v.elements, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (v Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(v.elements, tracker)

}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (v Vector[T]) For(walker func(T) error) error {
	return slice.For(v.elements, walker)
}

// ForEach applies the 'walker' function for every element
func (v Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(v.elements, walker)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (v Vector[T]) Filter(filter func(T) bool) stream.Iter[T] {
	h := v.Head()
	return stream.New(loop.Filter(h.Next, filter).Next)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (v Vector[T]) Filt(predicate func(T) (bool, error)) breakStream.Iter[T] {
	h := v.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (v Vector[T]) Convert(converter func(T) T) stream.Iter[T] {
	return collection.Convert(v, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (v Vector[T]) Conv(converter func(T) (T, error)) breakStream.Iter[T] {
	return collection.Conv(v, converter)
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
