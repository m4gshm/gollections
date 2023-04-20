package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys instantiates MapKeys using the elements as internal storage
func WrapKeys[K comparable, V any](elements map[K]V) *MapKeys[K, V] {
	return &MapKeys[K, V]{elements}
}

// MapKeys is the container reveal keys of a map and hides values
type MapKeys[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Collection[int] = (*MapKeys[int, any])(nil)
	_ fmt.Stringer      = (*MapValues[int, any])(nil)
)

// Begin creates iterator
func (m *MapKeys[K, V]) Begin() c.Iterator[K] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m *MapKeys[K, V]) Head() iter.Key[K, V] {
	var elements map[K]V
	if m != nil {
		elements = m.elements
	}
	return *iter.NewKey(elements)
}

// First returns the first element of the collection, an iterator to iterate over the remaining elements, and true\false marker of availability next elements.
// If no more elements then ok==false.
func (m *MapKeys[K, V]) First() (iter.Key[K, V], K, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

// Len returns amount of elements
func (m *MapKeys[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.elements)
}

// IsEmpty returns true if the collection is empty
func (m *MapKeys[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Slice collects the elements to a slice
func (m *MapKeys[K, V]) Slice() (keys []K) {
	if m != nil {
		keys = map_.Keys(m.elements)
	}
	return keys
}

// For applies the 'walker' function for every key. Return the c.ErrBreak to stop.
func (m *MapKeys[K, V]) For(walker func(K) error) error {
	if m == nil {
		return nil
	}
	return map_.ForKeys(m.elements, walker)
}

// ForEach applies the 'walker' function for every key
func (m *MapKeys[K, V]) ForEach(walker func(K)) {
	if m != nil {
		map_.ForEachKey(m.elements, walker)
	}
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (m *MapKeys[K, V]) Filter(filter func(K) bool) c.Pipe[K] {
	h := m.Head()
	return iter.NewPipe[K](iter.Filter(h, h.Next, filter))
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (m *MapKeys[K, V]) Convert(converter func(K) K) c.Pipe[K] {
	h := m.Head()
	return iter.NewPipe[K](iter.Convert(h, h.Next, converter))
}

// Reduce reduces the elements into an one using the 'merge' function
func (m *MapKeys[K, V]) Reduce(by func(K, K) K) K {
	h := m.Head()
	return loop.Reduce(h.Next, by)
}

func (m *MapKeys[K, V]) String() string {
	return slice.ToString(m.Slice())
}
