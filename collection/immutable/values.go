package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice"
)

// WrapVal instantiates MapValues using elements as internal storage.
func WrapVal[K comparable, V any](elements map[K]V) MapValues[K, V] {
	return MapValues[K, V]{elements}
}

// MapValues is the wrapper for Map'm values.
type MapValues[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ collection.Collection[any] = (*MapValues[int, any])(nil)
	_ collection.Collection[any] = MapValues[int, any]{}
	_ fmt.Stringer               = (*MapValues[int, any])(nil)
	_ fmt.Stringer               = MapValues[int, any]{}
)

// All is used to iterate through the collection using `for val := range`.
func (m MapValues[K, V]) All(consumer func(V) bool) {
	map_.TrackValuesWhile(m.elements, consumer)
}

// Head returns the first element.
func (m MapValues[K, V]) Head() (V, bool) {
	return collection.Head(m)
}

// Len returns amount of elements
func (m MapValues[K, V]) Len() int {
	return len(m.elements)
}

// IsEmpty returns true if the collection is empty
func (m MapValues[K, V]) IsEmpty() bool {
	return collection.IsEmpty(m)
}

// Slice collects the values to a slice
func (m MapValues[K, V]) Slice() []V {
	return map_.Values(m.elements)
}

// Append collects the values to the specified 'out' slice
func (m MapValues[K, V]) Append(out []V) []V {
	return map_.AppendValues(m.elements, out)
}

// ForEach applies the 'consumer' function for every value of the collection
func (m MapValues[K, V]) ForEach(consumer func(V)) {
	map_.ForEachValue(m.elements, consumer)
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filter(filter func(V) bool) seq.Seq[V] {
	return collection.Filter(m, filter)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filt(predicate func(V) (bool, error)) seq.SeqE[V] {
	return collection.Filt(m, predicate)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Convert(converter func(V) V) seq.Seq[V] {
	return collection.Convert(m, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Conv(converter func(V) (V, error)) seq.SeqE[V] {
	return collection.Conv(m, converter)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapValues[K, V]) Reduce(merge func(V, V) V) V {
	_, v := map_.Reduce(m.elements, func(_, _ K, v1, v2 V) (k K, v V) {
		return k, merge(v1, v2)
	})
	return v
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (m MapValues[K, V]) HasAny(predicate func(V) bool) bool {
	return map_.HasAny(m.elements, func(_ K, v V) bool {
		return predicate(v)
	})
}

// Sort creates a vector with sorted the values
func (m MapValues[K, V]) Sort(comparer slice.Comparer[V]) Vector[V] {
	return WrapVector(slice.Sort(m.Slice(), comparer))
}

func (m MapValues[K, V]) String() string {
	return slice.ToString(m.Slice())
}
