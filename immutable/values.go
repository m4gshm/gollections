package immutable

import (
	"fmt"
	"sort"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
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
	_ c.Collection[any, *map_.NewValIter[int, any]] = (*MapValues[int, any])(nil)
	_ fmt.Stringer                                  = (*MapValues[int, any])(nil)

	_ c.Collection[any, *map_.NewValIter[int, any]] = MapValues[int, any]{}
	_ fmt.Stringer                                  = MapValues[int, any]{}
)

// Begin creates iterator
func (m MapValues[K, V]) Begin() *map_.NewValIter[K, V] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m MapValues[K, V]) Head() map_.NewValIter[K, V] {
	return map_.NewVal(m.elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapValues[K, V]) First() (map_.NewValIter[K, V], V, bool) {
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
	return stream.New(loop.Filter(h.Next, filter).Next)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filt(filter func(V) (bool, error)) breakLoop.StreamIter[V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Filt(breakLoop.From(h.Next), filter).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Convert(converter func(V) V) stream.Iter[V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, converter).Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m MapValues[K, V]) Conv(converter func(V) (V, error)) breakLoop.StreamIter[V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Conv(breakLoop.From(h.Next), converter).Next)
}

// Reduce reduces the elements into an one using the 'merge' function
func (m MapValues[K, V]) Reduce(merge func(V, V) V) V {
	_, v := map_.Reduce(m.elements, func(_ K, v1 V, _ K, v2 V) (k K, v V) {
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
func (m MapValues[K, V]) Sort(less func(e1, e2 V) bool) Vector[V] {
	var dest = m.Slice()
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (m MapValues[K, V]) String() string {
	return slice.ToString(m.Slice())
}
