package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapVal instantiates MapValues using elements as internal storage.
func WrapVal[K comparable, V any](order []K, elements map[K]V) MapValues[K, V] {
	return MapValues[K, V]{order, elements}
}

// MapValues is the wrapper for Map's values.
type MapValues[K comparable, V any] struct {
	order    []K
	elements map[K]V
}

var (
	_ c.Collection[any] = (*MapValues[int, any])(nil)
	_ c.Collection[any] = MapValues[int, any]{}
	_ fmt.Stringer      = (*MapValues[int, any])(nil)
	_ fmt.Stringer      = MapValues[int, any]{}
)

func (s MapValues[K, V]) Begin() c.Iterator[V] {
	return s.Head()
}

func (s MapValues[K, V]) Head() *ValIter[K, V] {
	return NewValIter(s.order, s.elements)
}

func (s MapValues[K, V]) First() (*ValIter[K, V], V, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (s MapValues[K, V]) Len() int {
	return len(s.elements)
}

func (s MapValues[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s MapValues[K, V]) Slice() []V {
	elements := make([]V, len(s.order))
	for i, key := range s.order {
		val := s.elements[key]
		elements[i] = val
	}
	return elements
}

func (s MapValues[K, V]) For(walker func(V) error) error {
	return map_.ForOrderedValues(s.order, s.elements, walker)
}

func (s MapValues[K, V]) ForEach(walker func(V)) {
	map_.ForEachOrderedValues(s.order, s.elements, walker)
}

func (s MapValues[K, V]) Get(index int) (V, bool) {
	keys := s.order
	if index >= 0 && index < len(keys) {
		key := keys[index]
		val, ok := s.elements[key]
		return val, ok
	}
	var no V
	return no, false
}

func (s MapValues[K, V]) Filter(filter func(V) bool) c.Pipe[V] {
	h := s.Head()
	return iter.NewPipe[V](iter.Filter(h, h.Next, filter))
}

func (s MapValues[K, V]) Convert(by func(V) V) c.Pipe[V] {
	h := s.Head()
	return iter.NewPipe[V](iter.Convert(h, h.Next, by))
}

func (s MapValues[K, V]) Reduce(by func(V, V) V) V {
	return loop.Reduce(s.Head().Next, by)
}

func (s MapValues[K, V]) String() string {
	return slice.ToString(s.Slice())
}
