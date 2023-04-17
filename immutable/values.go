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
func WrapVal[K comparable, V any](elements map[K]V) MapValues[K, V] {
	return MapValues[K, V]{elements}
}

// MapValues is the wrapper for Map's values.
type MapValues[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Collection[any] = (*MapValues[int, any])(nil)
	_ c.Collection[any] = MapValues[int, any]{}
	_ fmt.Stringer      = (*MapValues[int, any])(nil)
	_ fmt.Stringer      = MapValues[int, any]{}
)

func (s MapValues[K, V]) Begin() c.Iterator[V] {
	h := s.Head()
	return &h
}

func (s MapValues[K, V]) Head() iter.Val[K, V] {
	return iter.NewVal(s.elements)
}

func (s MapValues[K, V]) First() (iter.Val[K, V], V, bool) {
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
	elements := make([]V, 0, len(s.elements))
	for _, val := range s.elements {
		elements = append(elements, val)
	}
	return elements
}

func (s MapValues[K, V]) For(walker func(V) error) error {
	return map_.ForValues(s.elements, walker)
}

func (s MapValues[K, V]) ForEach(walker func(V)) {
	map_.ForEachValue(s.elements, walker)
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
	h := s.Head()
	return loop.Reduce(h.Next, by)
}

func (s MapValues[K, V]) Sort(less func(e1, e2 V) bool) Vector[V] {
	var dest = s.Slice()
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (s MapValues[K, V]) String() string {
	return slice.ToString(s.Slice())
}
