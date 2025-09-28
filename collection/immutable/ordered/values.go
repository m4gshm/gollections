package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapVal instantiates MapValues using elements as internal storage.
func WrapVal[K comparable, V any](order []K, elements map[K]V) MapValues[K, V] {
	return MapValues[K, V]{order, elements}
}

// MapValues is the wrapper for Map'm values.
type MapValues[K comparable, V any] struct {
	order    []K
	elements map[K]V
}

var (
	_ collection.Collection[any] = (*MapValues[int, any])(nil)
	_ c.OrderedRange[any]        = (*MapValues[int, any])(nil)
	_ fmt.Stringer               = (*MapValues[int, any])(nil)
)

// Head creates an iterator to iterate through the collection.
//
// Deprecated: replaced by [MapValues.All].
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
func (m MapValues[K, V]) Slice() (values []V) {
	return m.Append(values)
}

// Append collects the values to the specified 'out' slice
func (m MapValues[K, V]) Append(out []V) (values []V) {
	for _, key := range m.order {
		val := m.elements[key]
		out = append(out, val)
	}
	return out
}

// All is used to iterate through the collection using `for val := range`.
func (m MapValues[K, V]) All(consumer func(V) bool) {
	m.IAll(func(_ int, v V) bool { return consumer(v) })
}

// IAll is used to iterate through the collection using `for index, val := range`.
func (m MapValues[K, V]) IAll(consumer func(int, V) bool) {
	map_.TrackOrderedValuesWhile(m.order, m.elements, consumer)
}

// For applies the 'consumer' function for every value until the consumer returns the c.Break to stop.
func (m MapValues[K, V]) For(consumer func(V) error) error {
	return map_.ForOrderedValues(m.order, m.elements, consumer)
}

// ForEach applies the 'consumer' function for every value
func (m MapValues[K, V]) ForEach(consumer func(V)) {
	map_.ForEachOrderedValues(m.order, m.elements, consumer)
}

// Get returns an element by the index, otherwise, if the provided index is ouf of the collection len, returns zero T and false in the second result
func (m MapValues[K, V]) Get(index int) (V, bool) {
	keys := m.order
	if index >= 0 && index < len(keys) {
		key := keys[index]
		val, ok := m.elements[key]
		return val, ok
	}
	var no V
	return no, false
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filter(filter func(V) bool) collection.Seq[V] {
	return collection.Filter(m, filter)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filt(filter func(V) (bool, error)) collection.SeqE[V] {
	return collection.Filt(m, filter)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Convert(converter func(V) V) collection.Seq[V] {
	return collection.Convert(m, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Conv(converter func(V) (V, error)) collection.SeqE[V] {
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

func (m MapValues[K, V]) String() string {
	return slice.ToString(m.Slice())
}
