package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
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
	_ collection.Collection[int] = (*MapKeys[int])(nil)
	_ collection.Collection[int] = MapKeys[int]{}
	_ fmt.Stringer               = (*MapKeys[int])(nil)
	_ fmt.Stringer               = MapKeys[int]{}
)

// Loop creates a loop to iterating through elements.
func (m MapKeys[K]) Loop() loop.Loop[K] {
	return loop.Of(m.keys...)
}

// Head creates an iterator and returns as implementation type value
func (m MapKeys[K]) Head() slice.Iter[K] {
	return slice.NewHead(m.keys)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapKeys[K]) First() (slice.Iter[K], K, bool) {
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
func (m MapKeys[K]) Slice() (out []K) {
	if keys := m.keys; keys != nil {
		out = slice.Clone(keys)
	}
	return out
}

// Append collects the values to the specified 'out' slice
func (m MapKeys[K]) Append(out []K) []K {
	if keys := m.keys; keys != nil {
		out = append(out, keys...)
	}
	return out
}

// For applies the 'walker' function for every key. Return the c.ErrBreak to stop.
func (m MapKeys[K]) For(walker func(K) error) error {
	return slice.For(m.keys, walker)
}

// ForEach applies the 'walker' function for every element
func (m MapKeys[K]) ForEach(walker func(K)) {
	slice.ForEach(m.keys, walker)
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K]) Filter(filter func(K) bool) loop.Loop[K] {
	return loop.Filter(m.Loop(), filter)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K]) Filt(predicate func(K) (bool, error)) breakLoop.Loop[K] {
	return loop.Filt(m.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (m MapKeys[K]) Convert(converter func(K) K) loop.Loop[K] {
	return loop.Convert(m.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (m MapKeys[K]) Conv(converter func(K) (K, error)) breakLoop.Loop[K] {
	return loop.Conv(m.Loop(), converter)
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
	return slice.Gett(m.keys, index)
}
