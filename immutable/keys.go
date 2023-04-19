package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys is non-copy constructor
func WrapKeys[K comparable, V any](uniques map[K]V) *MapKeys[K, V] {
	return &MapKeys[K, V]{uniques}
}

// MapKeys is the container reveal keys of a map and hides values.
type MapKeys[K comparable, V any] struct {
	uniques map[K]V
}

var (
	_ c.Collection[int] = (*MapKeys[int, any])(nil)
	_ fmt.Stringer      = (*MapValues[int, any])(nil)
)

func (m *MapKeys[K, V]) Begin() c.Iterator[K] {
	h := m.Head()
	return &h
}

func (m *MapKeys[K, V]) Head() iter.Key[K, V] {
	var uniques map[K]V
	if m != nil {
		uniques = m.uniques
	}
	return *iter.NewKey(uniques)
}

func (m *MapKeys[K, V]) First() (iter.Key[K, V], K, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (m *MapKeys[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.uniques)
}

func (m *MapKeys[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *MapKeys[K, V]) Slice() (elements []K) {
	if m != nil {
		uniques := m.uniques
		elements = make([]K, 0, len(uniques))
		for key := range uniques {
			elements = append(elements, key)
		}
	}
	return elements
}

func (m *MapKeys[K, V]) For(walker func(K) error) error {
	if m == nil {
		return nil
	}
	return map_.ForKeys(m.uniques, walker)
}

func (m *MapKeys[K, V]) ForEach(walker func(K)) {
	if m != nil {
		map_.ForEachKey(m.uniques, walker)
	}
}

func (m *MapKeys[K, V]) Filter(filter func(K) bool) c.Pipe[K] {
	h := m.Head()
	return iter.NewPipe[K](iter.Filter(h, h.Next, filter))
}

func (m *MapKeys[K, V]) Convert(by func(K) K) c.Pipe[K] {
	h := m.Head()
	return iter.NewPipe[K](iter.Convert(h, h.Next, by))
}

func (m *MapKeys[K, V]) Reduce(by func(K, K) K) K {
	h := m.Head()
	return loop.Reduce((&h).Next, by)
}

func (m *MapKeys[K, V]) String() string {
	return slice.ToString(m.Slice())
}
