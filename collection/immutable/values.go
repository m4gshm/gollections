package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
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

// Loop creates a loop to iterating through elements.
func (m MapValues[K, V]) Loop() loop.Loop[V] {
	h := m.Head()
	return (&h).Next
}

// Head creates an iterator and returns as implementation type value
func (m MapValues[K, V]) Head() map_.ValIter[K, V] {
	return map_.NewValIter(m.elements)
}

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
	return m.Len() == 0
}

// Slice collects the values to a slice
func (m MapValues[K, V]) Slice() []V {
	return map_.Values(m.elements)
}

// Append collects the values to the specified 'out' slice
func (m MapValues[K, V]) Append(out []V) []V {
	return map_.AppendValues(m.elements, out)
}

// For applies the 'walker' function for collection values. Return the c.ErrBreak to stop.
func (m MapValues[K, V]) For(walker func(V) error) error {
	return map_.ForValues(m.elements, walker)
}

// ForEach applies the 'walker' function for every value of the collection
func (m MapValues[K, V]) ForEach(walker func(V)) {
	map_.ForEachValue(m.elements, walker)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filter(filter func(V) bool) stream.Iter[V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter))
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filt(predicate func(V) (bool, error)) breakStream.Iter[V] {
	h := m.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate))
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Convert(converter func(V) V) stream.Iter[V] {
	return collection.Convert(m, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Conv(converter func(V) (V, error)) breakStream.Iter[V] {
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
