package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/seq"
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
	_ c.OrderedRange[any]    = Vector[any]{}
	_ fmt.Stringer           = (*Vector[any])(nil)
	_ fmt.Stringer           = Vector[any]{}
)

// All is used to iterate through the collection using `for e := range`.
func (v Vector[T]) All(consumer func(T) bool) {
	slice.WalkWhile(v.elements, consumer)
}

// IAll is used to iterate through the collection using `for i, e := range`.
func (v Vector[T]) IAll(consumer func(int, T) bool) {
	slice.TrackWhile(v.elements, consumer)
}

// Head returns the first element.
func (v Vector[T]) Head() (T, bool) {
	return collection.Head(v)
}

// Tail returns the latest element.
func (v Vector[T]) Tail() (T, bool) {
	return slice.Tail(v.elements)
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
	return collection.IsEmpty(v)
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the vector len, returns zero T and false in the second result
func (v Vector[T]) Get(index int) (out T, ok bool) {
	return slice.Gett(v.elements, index)
}

// TrackEach applies the 'consumer' function for every key/value pairs
func (v Vector[T]) TrackEach(consumer func(int, T)) {
	slice.TrackEach(v.elements, consumer)
}

// ForEach applies the 'consumer' function for every element
func (v Vector[T]) ForEach(consumer func(T)) {
	slice.ForEach(v.elements, consumer)
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'filter' function
func (v Vector[T]) Filter(filter func(T) bool) seq.Seq[T] {
	return collection.Filter(v, filter)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'filter' function
func (v Vector[T]) Filt(filter func(T) (bool, error)) seq.SeqE[T] {
	return collection.Filt(v, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (v Vector[T]) Convert(converter func(T) T) seq.Seq[T] {
	return collection.Convert(v, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (v Vector[T]) Conv(converter func(T) (T, error)) seq.SeqE[T] {
	return collection.Conv(v, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (v Vector[T]) Reduce(merge func(T, T) T) T {
	return slice.Reduce(v.elements, merge)
}

// HasAny checks whether the vector contains an element that satisfies the condition.
func (v Vector[T]) HasAny(condition func(T) bool) bool {
	return slice.HasAny(v.elements, condition)
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
