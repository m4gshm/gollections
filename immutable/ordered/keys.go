package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakIter "github.com/m4gshm/gollections/break/slice/iter"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/iter"
)

// WrapKeys instantiates MapKeys using the elements as internal storage
func WrapKeys[K comparable](elements []K) MapKeysIter[K] {
	return MapKeysIter[K]{elements}
}

// MapKeysIter is the wrapper for Map'm keys
type MapKeysIter[K comparable] struct {
	keys []K
}

var (
	_ c.Collection[int, *slice.Iter[int]] = (*MapKeysIter[int])(nil)
	_ fmt.Stringer                        = (*MapKeysIter[int])(nil)

	_ c.Collection[int, *slice.Iter[int]] = MapKeysIter[int]{}
	_ fmt.Stringer                        = MapKeysIter[int]{}
)

// Begin creates iterator
func (m MapKeysIter[K]) Begin() *slice.Iter[K] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m MapKeysIter[K]) Head() slice.Iter[K] {
	return slice.NewHead(m.keys)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapKeysIter[K]) First() (slice.Iter[K], K, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Len returns amount of elements
func (m MapKeysIter[K]) Len() int {
	return len(m.keys)
}

// IsEmpty returns true if the collection is empty
func (m MapKeysIter[K]) IsEmpty() bool {
	return m.Len() == 0
}

// Slice collects the elements to a slice
func (m MapKeysIter[K]) Slice() []K {
	return slice.Clone(m.keys)
}

// For applies the 'walker' function for every key. Return the c.ErrBreak to stop.
func (m MapKeysIter[K]) For(walker func(K) error) error {
	return slice.For(m.keys, walker)
}

// ForEach applies the 'walker' function for every element
func (m MapKeysIter[K]) ForEach(walker func(K)) {
	slice.ForEach(m.keys, walker)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeysIter[K]) Filter(filter func(K) bool) loop.StreamIter[K] {
	f := iter.Filter(m.keys, filter)
	return loop.Stream(f.Next)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeysIter[K]) Filt(filter func(K) (bool, error)) breakLoop.StreamIter[K] {
	f := breakIter.Filt(m.keys, filter)
	return breakLoop.Stream(f.Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m MapKeysIter[K]) Convert(converter func(K) K) loop.StreamIter[K] {
	conv := iter.Convert(m.keys, converter)
	return loop.Stream(conv.Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m MapKeysIter[K]) Conv(converter func(K) (K, error)) breakLoop.StreamIter[K] {
	conv := breakIter.Conv(m.keys, converter)
	return breakLoop.Stream(conv.Next)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapKeysIter[K]) Reduce(merge func(K, K) K) K {
	return slice.Reduce(m.keys, merge)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (m MapKeysIter[K]) HasAny(predicate func(K) bool) bool {
	return slice.HasAny(m.keys, predicate)
}

// String returns string representation of the collection
func (m MapKeysIter[K]) String() string {
	return slice.ToString(m.Slice())
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the collection len, returns zero T and false in the second result
func (m MapKeysIter[K]) Get(index int) (K, bool) {
	return slice.Gett(m.keys, index)
}
