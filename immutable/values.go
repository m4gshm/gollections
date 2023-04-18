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

func (m *MapValues[K, V]) Begin() c.Iterator[V] {
	h := m.Head()
	return &h
}

func (m *MapValues[K, V]) Head() iter.Val[K, V] {
	var elements map[K]V
	if m != nil {
		elements = m.elements
	}
	return *iter.NewVal(elements)
}

func (m *MapValues[K, V]) First() (iter.Val[K, V], V, bool) {
	var (
		iterator  = m.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (m *MapValues[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.elements)
}

func (m *MapValues[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *MapValues[K, V]) Slice() (elements []V) {
	if m == nil {
		return
	}
	elements = make([]V, 0, len(m.elements))
	for _, val := range m.elements {
		elements = append(elements, val)
	}
	return elements
}

func (m *MapValues[K, V]) For(walker func(V) error) error {
	if m == nil {
		return nil
	}
	return map_.ForValues(m.elements, walker)
}

func (m *MapValues[K, V]) ForEach(walker func(V)) {
	if m != nil {
		map_.ForEachValue(m.elements, walker)
	}
}

func (m *MapValues[K, V]) Filter(filter func(V) bool) c.Pipe[V] {
	h := m.Head()
	return iter.NewPipe[V](iter.Filter(h, h.Next, filter))
}

func (m *MapValues[K, V]) Convert(by func(V) V) c.Pipe[V] {
	h := m.Head()
	return iter.NewPipe[V](iter.Convert(h, h.Next, by))
}

func (m *MapValues[K, V]) Reduce(by func(V, V) V) V {
	h := m.Head()
	return loop.Reduce(h.Next, by)
}

func (m *MapValues[K, V]) Sort(less func(e1, e2 V) bool) *Vector[V] {
	var dest = m.Slice()
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (m *MapValues[K, V]) String() string {
	return slice.ToString(m.Slice())
}
