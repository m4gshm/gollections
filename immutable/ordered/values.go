package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	omap "github.com/m4gshm/gollections/immutable/ordered/map_"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/stream"
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
	_ c.Collection[any]                         = (*MapValues[int, any])(nil)
	_ loop.Looper[any, *omap.ValIter[int, any]] = (*MapValues[int, any])(nil)
	_ fmt.Stringer                              = (*MapValues[int, any])(nil)
)

// Iter creates an iterator and returns as interface
func (m MapValues[K, V]) Iter() c.Iterator[V] {
	h := m.Head()
	return &h
}

// Loop creates an iterator and returns as implementation type reference
func (m MapValues[K, V]) Loop() *omap.ValIter[K, V] {
	h := m.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (m MapValues[K, V]) Head() omap.ValIter[K, V] {
	var (
		order    []K
		elements map[K]V
	)

	order = m.order
	elements = m.elements

	return omap.NewValIter(order, elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m MapValues[K, V]) First() (omap.ValIter[K, V], V, bool) {
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

// For applies the 'walker' function for every value. Return the c.ErrBreak to stop.
func (m MapValues[K, V]) For(walker func(V) error) error {
	return map_.ForOrderedValues(m.order, m.elements, walker)
}

// ForEach applies the 'walker' function for every value
func (m MapValues[K, V]) ForEach(walker func(V)) {
	map_.ForEachOrderedValues(m.order, m.elements, walker)
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

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filter(filter func(V) bool) stream.Iter[V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter).Next)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m MapValues[K, V]) Filt(filter func(V) (bool, error)) breakStream.Iter[V] {
	h := m.Head()
	return breakStream.New(breakLoop.Filt(breakLoop.From(h.Next), filter).Next)
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

func (m MapValues[K, V]) String() string {
	return slice.ToString(m.Slice())
}
