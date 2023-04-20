package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapVal instantiates MapValues using elements as internal storage.
func WrapVal[K comparable, V any](elements map[K]V) *MapValues[K, V] {
	return &MapValues[K, V]{elements}
}

// MapValues is the wrapper for Map'm values.
type MapValues[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Collection[any] = (*MapValues[int, any])(nil)
	_ fmt.Stringer      = (*MapValues[int, any])(nil)
)

// Begin creates iterator
func (m *MapValues[K, V]) Begin() c.Iterator[V] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m *MapValues[K, V]) Head() iter.Val[K, V] {
	var elements map[K]V
	if m != nil {
		elements = m.elements
	}
	return *iter.NewVal(elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m *MapValues[K, V]) First() (iter.Val[K, V], V, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Len returns amount of elements
func (m *MapValues[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.elements)
}

// IsEmpty returns true if the collection is empty
func (m *MapValues[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Slice collects the values to a slice
func (m *MapValues[K, V]) Slice() (values []V) {
	if m != nil {
		values = map_.Values(m.elements)

	}
	return values
}

// For applies the 'walker' function for collection values. Return the c.ErrBreak to stop.
func (m *MapValues[K, V]) For(walker func(V) error) error {
	if m == nil {
		return nil
	}
	return map_.ForValues(m.elements, walker)
}

// ForEach applies the 'walker' function for every value of the collection
func (m *MapValues[K, V]) ForEach(walker func(V)) {
	if m != nil {
		map_.ForEachValue(m.elements, walker)
	}
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (m *MapValues[K, V]) Filter(filter func(V) bool) c.Pipe[V] {
	h := m.Head()
	return iter.NewPipe[V](iter.Filter(h, h.Next, filter))
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (m *MapValues[K, V]) Convert(converter func(V) V) c.Pipe[V] {
	h := m.Head()
	return iter.NewPipe[V](iter.Convert(h, h.Next, converter))
}

// Reduce reduces the elements into an one using the 'merge' function
func (m *MapValues[K, V]) Reduce(merge func(V, V) V) V {
	h := m.Head()
	return loop.Reduce(h.Next, merge)
}

// Sort creates a vector with sorted the values
func (m *MapValues[K, V]) Sort(less func(e1, e2 V) bool) *Vector[V] {
	var dest = m.Slice()
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (m *MapValues[K, V]) String() string {
	return slice.ToString(m.Slice())
}
