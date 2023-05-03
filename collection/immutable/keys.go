package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
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
	_ c.Collection[int]                         = (*MapKeys[int, any])(nil)
	_ c.Collection[int]                         = MapKeys[int, any]{}
	_ loop.Looper[int, *map_.KeyIter[int, any]] = (*MapKeys[int, any])(nil)
	_ loop.Looper[int, *map_.KeyIter[int, any]] = MapKeys[int, any]{}
	_ fmt.Stringer                              = (*MapKeys[int, any])(nil)
	_ fmt.Stringer                              = MapKeys[int, any]{}
)

// Iter creates an iterator and returns as interface
func (m MapKeys[K, V]) Iter() c.Iterator[K] {
	h := m.Head()
	return &h
}

// Loop creates an iterator and returns as implementation type reference
func (m MapKeys[K, V]) Loop() *map_.KeyIter[K, V] {
	h := m.Head()
	return &h
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

// For applies the 'walker' function for every key. Return the c.ErrBreak to stop.
func (m MapKeys[K, V]) For(walker func(K) error) error {
	return map_.ForKeys(m.elements, walker)
}

// ForEach applies the 'walker' function for every key
func (m MapKeys[K, V]) ForEach(walker func(K)) {
	map_.ForEachKey(m.elements, walker)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K, V]) Filter(filter func(K) bool) stream.Iter[K] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter).Next)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapKeys[K, V]) Filt(predicate func(K) (bool, error)) breakStream.Iter[K] {
	h := m.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m MapKeys[K, V]) Convert(converter func(K) K) stream.Iter[K] {
	return collection.Convert(m, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (m MapKeys[K, V]) Conv(converter func(K) (K, error)) breakStream.Iter[K] {
	return collection.Conv(m, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapKeys[K, V]) Reduce(merge func(K, K) K) K {
	k, _ := map_.Reduce(m.elements, func(k1 K, _ V, k2 K, _ V) (rk K, rv V) {
		return merge(k1, k2), rv
	})
	return k
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (m MapKeys[K, V]) HasAny(predicate func(K) bool) bool {
	return map_.HasAny(m.elements, func(k K, v V) bool {
		return predicate(k)
	})
}

func (m MapKeys[K, V]) String() string {
	return slice.ToString(m.Slice())
}
