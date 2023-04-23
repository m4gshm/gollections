package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	loopIter "github.com/m4gshm/gollections/loop/iter"
	"github.com/m4gshm/gollections/loop/stream"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/iter"
)

// NewVector instantiates Vector and copies elements to it.
func NewVector[T any](elements []T) Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// WrapVector instantiates Vector using a slise as internal storage.
func WrapVector[T any](elements []T) Vector[T] {
	return Vector[T]{elements: elements, esize: notsafe.GetTypeSize[T]()}
}

// Vector is the Collection implementation that provides elements order and index access.
type Vector[T any] struct {
	elements []T
	esize    uintptr
}

var (
	_ c.Vector[any] = (*Vector[any])(nil)
	_ fmt.Stringer  = (*Vector[any])(nil)
	_ c.Vector[any] = Vector[any]{}
	_ fmt.Stringer  = Vector[any]{}
)

// Begin creates iterator
func (v Vector[T]) Begin() c.Iterator[T] {
	h := v.Head()
	return &h
}

// Head creates iterator
func (v Vector[T]) Head() iter.ArrayIter[T] {
	return iter.NewHeadS(v.elements, v.esize)
}

// Tail creates an iterator pointing to the end of the collection
func (v Vector[T]) Tail() iter.ArrayIter[T] {
	return iter.NewTailS(v.elements, v.esize)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (v Vector[T]) First() (iter.ArrayIter[T], T, bool) {
	var (
		iterator  = iter.NewHeadS(v.elements, v.esize)
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Last returns the latest element of the collection, an iterator to reverse iterate over the remaining elements, and true\false marker of availability previous elements.
// If no more elements then ok==false.
func (v Vector[T]) Last() (iter.ArrayIter[T], T, bool) {
	var (
		iterator  = iter.NewTailS(v.elements, v.esize)
		first, ok = iterator.Prev()
	)
	return iterator, first, ok
}

// Slice collects the elements to a slice
func (v Vector[T]) Slice() []T {
	return slice.Clone(v.elements)
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
	return slice.Get(v.elements, index)
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

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (v Vector[T]) Filter(filter func(T) bool) c.Stream[T] {
	h := v.Head()
	return stream.New(loopIter.Filter(h.Next, filter).Next)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (v Vector[T]) Convert(converter func(T) T) c.Stream[T] {
	h := v.Head()
	return stream.New(loopIter.Convert(h.Next, converter).Next)
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
func (v Vector[T]) Sort(less slice.Less[T]) Vector[T] {
	return v.sortBy(sort.Slice, less)
}

// StableSort returns a stable sorted clone of the Vector
func (v Vector[T]) StableSort(less slice.Less[T]) Vector[T] {
	return v.sortBy(sort.SliceStable, less)
}

func (v Vector[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) Vector[T] {
	elements := slice.Clone(v.elements)
	slice.Sort(elements, sorter, less)
	return WrapVector(elements)
}

func (v Vector[T]) String() string {
	return slice.ToString(v.elements)
}
