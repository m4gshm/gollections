package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	iterslice "github.com/m4gshm/gollections/iter/impl/slice"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys instantiates MapKeys using the elements as internal storage
func WrapKeys[K comparable](elements []K) MapKeys[K] {
	return MapKeys[K]{elements}
}

// MapKeys is the wrapper for Map'm keys
type MapKeys[K comparable] struct {
	keys []K
}

var (
	_ c.Collection[int] = (*MapKeys[int])(nil)
	_ fmt.Stringer      = (*MapKeys[int])(nil)

	_ c.Collection[int] = MapKeys[int]{}
	_ fmt.Stringer      = MapKeys[int]{}
)

// Begin creates iterator
func (m MapKeys[K]) Begin() c.Iterator[K] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m MapKeys[K]) Head() iter.ArrayIter[K] {
	return iter.NewHead(m.keys)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapKeys[K]) First() (iter.ArrayIter[K], K, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Len returns amount of elements
func (m MapKeys[K]) Len() int {
	return len(m.keys)
}

// IsEmpty returns true if the collection is empty
func (m MapKeys[K]) IsEmpty() bool {
	return m.Len() == 0
}

// Slice collects the elements to a slice
func (m MapKeys[K]) Slice() []K {
	return slice.Clone(m.keys)
}

// For applies the 'walker' function for every key. Return the c.ErrBreak to stop.
func (m MapKeys[K]) For(walker func(K) error) error {
	return slice.For(m.keys, walker)
}

// ForEach applies the 'walker' function for every element
func (m MapKeys[K]) ForEach(walker func(K)) {
	slice.ForEach(m.keys, walker)
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K]) Filter(filter func(K) bool) c.Pipe[K] {
	f := iterslice.Filter(m.keys, filter)
	return iter.NewPipe[K](&f)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (m MapKeys[K]) Convert(converter func(K) K) c.Pipe[K] {
	conv := iterslice.Convert(m.keys, converter)
	return iter.NewPipe[K](&conv)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapKeys[K]) Reduce(merge func(K, K) K) K {
	return slice.Reduce(m.keys, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (m MapKeys[K]) HasAny(predicate func(K) bool) bool {
	return slice.HasAny(m.keys, predicate)
}

// String returns string representation of the collection
func (m MapKeys[K]) String() string {
	return slice.ToString(m.Slice())
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the collection len, returns zero T and false in the second result
func (m MapKeys[K]) Get(index int) (K, bool) {
	return slice.Get(m.keys, index)
}
