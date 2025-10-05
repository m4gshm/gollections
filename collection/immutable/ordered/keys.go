package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/seq"
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
	_ c.OrderedRange[int]        = MapKeys[int]{}
	_ fmt.Stringer               = (*MapKeys[int])(nil)
	_ fmt.Stringer               = MapKeys[int]{}
)

// All is used to iterate through the collection using `for key := range`.
func (m MapKeys[K]) All(consumer func(K) bool) {
	slice.WalkWhile(m.keys, consumer)
}

// IAll is used to iterate through the collection using `for index, key := range`.
func (m MapKeys[K]) IAll(consumer func(int, K) bool) {
	slice.TrackWhile(m.keys, consumer)
}

// Head returns the first element.
func (m MapKeys[K]) Head() (K, bool) {
	return collection.Head(m)
}

// Len returns amount of elements
func (m MapKeys[K]) Len() int {
	return len(m.keys)
}

// IsEmpty returns true if the collection is empty
func (m MapKeys[K]) IsEmpty() bool {
	return collection.IsEmpty(m)
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

// ForEach applies the 'consumer' function for every element
func (m MapKeys[K]) ForEach(consumer func(K)) {
	slice.ForEach(m.keys, consumer)
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'filter' function
func (m MapKeys[K]) Filter(filter func(K) bool) seq.Seq[K] {
	return collection.Filter(m, filter)
}

// Filt returns an errorable seq consisting of elements that satisfy the condition of the 'filter' function
func (m MapKeys[K]) Filt(filter func(K) (bool, error)) seq.SeqE[K] {
	return collection.Filt(m, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (m MapKeys[K]) Convert(converter func(K) K) seq.Seq[K] {
	return collection.Convert(m, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func (m MapKeys[K]) Conv(converter func(K) (K, error)) seq.SeqE[K] {
	return collection.Conv(m, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapKeys[K]) Reduce(merge func(K, K) K) K {
	return slice.Reduce(m.keys, merge)
}

// HasAny checks whether the collection contains a key that satisfies the condition.
func (m MapKeys[K]) HasAny(condition func(K) bool) bool {
	return slice.HasAny(m.keys, condition)
}

// First returns the first key that satisfies requirements of the condition
func (m MapKeys[K]) First(condition func(K) bool) (K, bool) {
	return slice.First(m.keys, condition)
}

// String returns string representation of the collection
func (m MapKeys[K]) String() string {
	return slice.ToString(m.Slice())
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the collection len, returns zero T and false in the second result
func (m MapKeys[K]) Get(index int) (K, bool) {
	return slice.Gett(m.keys, index)
}
