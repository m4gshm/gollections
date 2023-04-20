package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys instantiates MapKeys using the elements as internal storage
func WrapKeys[K comparable](elements []K) *MapKeys[K] {
	return &MapKeys[K]{elements}
}

// MapKeys is the wrapper for Map'm keys
type MapKeys[K comparable] struct {
	keys []K
}

var (
	_ c.Collection[int] = (*MapKeys[int])(nil)
	_ fmt.Stringer      = (*MapKeys[int])(nil)
)

// Begin creates iterator
func (m *MapKeys[K]) Begin() c.Iterator[K] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m *MapKeys[K]) Head() iter.ArrayIter[K] {
	var elements []K
	if m != nil {
		elements = m.keys
	}
	return iter.NewHead(elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m *MapKeys[K]) First() (iter.ArrayIter[K], K, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Len returns amount of elements
func (m *MapKeys[K]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.keys)
}

// IsEmpty returns true if the collection is empty
func (m *MapKeys[K]) IsEmpty() bool {
	return m.Len() == 0
}

// Slice collects the elements to a slice
func (m *MapKeys[K]) Slice() (keys []K) {
	if m != nil {
		keys = slice.Clone(m.keys)
	}
	return keys
}

// For applies the 'walker' function for every key. Return the c.ErrBreak to stop.
func (m *MapKeys[K]) For(walker func(K) error) error {
	if m == nil {
		return nil
	}
	return slice.For(m.keys, walker)
}

// ForEach applies the 'walker' function for every element
func (m *MapKeys[K]) ForEach(walker func(K)) {
	if m != nil {
		slice.ForEach(m.keys, walker)
	}
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (m *MapKeys[K]) Filter(filter func(K) bool) c.Pipe[K] {
	h := m.Head()
	return iter.NewPipe[K](iter.Filter(h, h.Next, filter))
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (m *MapKeys[K]) Convert(converter func(K) K) c.Pipe[K] {
	h := m.Head()
	return iter.NewPipe[K](iter.Convert(h, h.Next, converter))
}

// Reduce reduces the elements into an one using the 'merge' function
func (m *MapKeys[K]) Reduce(by func(K, K) K) K {
	h := m.Head()
	return loop.Reduce(h.Next, by)
}

// String returns string representation of the collection
func (m *MapKeys[K]) String() string {
	return slice.ToString(m.Slice())
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the collection len, returns zero T and false in the second result
func (m *MapKeys[K]) Get(index int) (K, bool) {
	return slice.Get(m.keys, index)
}
