package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys instantiates MapKeys using the elements as internal storage
func WrapKeys[K comparable, V any](elements map[K]V) MapKeys[K, V] {
	return MapKeys[K, V]{elements}
}

// MapKeys is the container reveal keys of a map and hides values
type MapKeys[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ collection.Collection[int] = (*MapKeys[int, any])(nil)
	_ collection.Collection[int] = MapKeys[int, any]{}
	_ fmt.Stringer               = (*MapKeys[int, any])(nil)
	_ fmt.Stringer               = MapKeys[int, any]{}
)

func (m MapKeys[K, V]) All(consumer func(K) bool) {
	map_.TrackKeysWhile(m.elements, consumer)
}

// Loop creates a loop to iterating through elements.
func (m MapKeys[K, V]) Loop() loop.Loop[K] {
	h := m.Head()
	return (&h).Next
}

// Head creates an iterator and returns as implementation type value
func (m MapKeys[K, V]) Head() map_.KeyIter[K, V] {
	return map_.NewKeyIter(m.elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapKeys[K, V]) First() (map_.KeyIter[K, V], K, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Len returns amount of elements
func (m MapKeys[K, V]) Len() int {
	return len(m.elements)
}

// IsEmpty returns true if the collection is empty
func (m MapKeys[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Slice collects the elements to a slice
func (m MapKeys[K, V]) Slice() []K {
	return map_.Keys(m.elements)
}

// Append collects the values to the specified 'out' slice
func (m MapKeys[K, V]) Append(out []K) []K {
	return map_.AppendKeys(m.elements, out)
}

// For applies the 'consumer' function for every key. Return the c.Break to stop.
func (m MapKeys[K, V]) For(consumer func(K) error) error {
	return map_.ForKeys(m.elements, consumer)
}

// ForEach applies the 'consumer' function for every key
func (m MapKeys[K, V]) ForEach(consumer func(K)) {
	map_.ForEachKey(m.elements, consumer)
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K, V]) Filter(filter func(K) bool) loop.Loop[K] {
	return loop.Filter(m.Loop(), filter)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K, V]) Filt(predicate func(K) (bool, error)) breakLoop.Loop[K] {
	return loop.Filt(m.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (m MapKeys[K, V]) Convert(converter func(K) K) loop.Loop[K] {
	return loop.Convert(m.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (m MapKeys[K, V]) Conv(converter func(K) (K, error)) breakLoop.Loop[K] {
	return loop.Conv(m.Loop(), converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapKeys[K, V]) Reduce(merge func(K, K) K) K {
	k, _ := map_.Reduce(m.elements, func(k1, k2 K, _, _ V) (rk K, rv V) {
		return merge(k1, k2), rv
	})
	return k
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (m MapKeys[K, V]) HasAny(predicate func(K) bool) bool {
	return map_.HasAny(m.elements, func(k K, _ V) bool {
		return predicate(k)
	})
}

func (m MapKeys[K, V]) String() string {
	return slice.ToString(m.Slice())
}
