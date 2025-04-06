package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
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

// All is used to iterate through the collection using `for ... range`. Supported since go 1.22 with GOEXPERIMENT=rangefunc enabled.
func (m MapValues[K, V]) All(consumer func(V) bool) {
	map_.TrackValuesWhile(m.elements, consumer)
}

// Loop creates a loop to iterate through the collection.
func (m MapValues[K, V]) Loop() loop.Loop[V] {
	h := m.Head()
	return (&h).Next
}

// Deprecated: Head is deprecated. Will be replaced by rance-over function iterator.
// Head creates an iterator to iterate through the collection.
func (m MapValues[K, V]) Head() map_.ValIter[K, V] {
	return map_.NewValIter(m.elements)
}

// Deprecated: First is deprecated. Will be replaced by rance-over function iterator.
// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapValues[K, V]) First() (map_.ValIter[K, V], V, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
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

// For applies the 'consumer' function for collection values until the consumer returns the c.Break to stop.
func (m MapValues[K, V]) For(consumer func(V) error) error {
	return map_.ForValues(m.elements, consumer)
}

// ForEach applies the 'consumer' function for every value of the collection
func (m MapValues[K, V]) ForEach(consumer func(V)) {
	map_.ForEachValue(m.elements, consumer)
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filter(filter func(V) bool) loop.Loop[V] {
	h := m.Head()
	return loop.Filter(h.Next, filter)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filt(predicate func(V) (bool, error)) breakLoop.Loop[V] {
	return loop.Filt(m.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Convert(converter func(V) V) loop.Loop[V] {
	return loop.Convert(m.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Conv(converter func(V) (V, error)) breakLoop.Loop[V] {
	return loop.Conv(m.Loop(), converter)
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
