package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/kv/predicate"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
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

// All is used to iterate through the collection using `for key := range`.
func (m MapKeys[K, V]) All(consumer func(K) bool) {
	map_.TrackKeysWhile(m.elements, consumer)
}

// Head returns the first element.
func (m MapKeys[K, V]) Head() (K, bool) {
	return collection.Head(m)
}

// Len returns amount of elements
func (m MapKeys[K, V]) Len() int {
	return len(m.elements)
}

// IsEmpty returns true if the collection is empty
func (m MapKeys[K, V]) IsEmpty() bool {
	return collection.IsEmpty(m)
}

// Slice collects the elements to a slice
func (m MapKeys[K, V]) Slice() []K {
	return map_.Keys(m.elements)
}

// Append collects the values to the specified 'out' slice
func (m MapKeys[K, V]) Append(out []K) []K {
	return map_.AppendKeys(m.elements, out)
}

// ForEach applies the 'consumer' function for every key
func (m MapKeys[K, V]) ForEach(consumer func(K)) {
	map_.ForEachKey(m.elements, consumer)
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'filter' function
func (m MapKeys[K, V]) Filter(filter func(K) bool) seq.Seq[K] {
	return collection.Filter(m, filter)
}

// Filt returns an errorable seq consisting of elements that satisfy the condition of the 'filter' function
func (m MapKeys[K, V]) Filt(filter func(K) (bool, error)) seq.SeqE[K] {
	return collection.Filt(m, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (m MapKeys[K, V]) Convert(converter func(K) K) seq.Seq[K] {
	return collection.Convert(m, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func (m MapKeys[K, V]) Conv(converter func(K) (K, error)) seq.SeqE[K] {
	return collection.Conv(m, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapKeys[K, V]) Reduce(merge func(K, K) K) K {
	k, _ := map_.Reduce(m.elements, func(k1, k2 K, _, _ V) (rk K, rv V) {
		return merge(k1, k2), rv
	})
	return k
}

// HasAny checks whether the collection contains a key that satisfies the condition.
func (m MapKeys[K, V]) HasAny(condition func(K) bool) bool {
	return map_.HasAny(m.elements, predicate.Key[V](condition))
}

// HasAny checks whether the collection contains a key that satisfies the condition.
func (m MapKeys[K, V]) First(condition func(K) bool) (K, bool) {
	k, _, ok := map_.First(m.elements, predicate.Key[V](condition))
	return k, ok
}

func (m MapKeys[K, V]) String() string {
	return slice.ToString(m.Slice())
}
