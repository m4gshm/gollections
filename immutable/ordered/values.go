package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func WrapVal[K comparable, V any](elements []K, uniques map[K]V) *MapValues[K, V] {
	return &MapValues[K, V]{elements, uniques}
}

type MapValues[K comparable, V any] struct {
	elements []K
	uniques  map[K]V
}

var _ c.Collection[any, []any, c.Iterator[any]] = (*MapValues[int, any])(nil)
var _ fmt.Stringer = (*MapValues[int, any])(nil)

func (s *MapValues[K, V]) Begin() c.Iterator[V] {
	return s.Head()
}

func (s *MapValues[K, V]) Head() *ValIter[K, V] {
	return NewValIter(s.elements, s.uniques)
}

func (s *MapValues[K, V]) Len() int {
	return len(s.uniques)
}

func (s *MapValues[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *MapValues[K, V]) Collect() []V {
	refs := s.elements
	elements := make([]V, len(refs))
	for i, key := range refs {
		val := s.uniques[key]
		elements[i] = val
	}
	return elements
}

func (s *MapValues[K, V]) For(walker func(V) error) error {
	return map_.ForOrderedValues(s.elements, s.uniques, walker)
}

func (s *MapValues[K, V]) ForEach(walker func(V)) {
	map_.ForEachOrderedValues(s.elements, s.uniques, walker)
}

func (s *MapValues[K, V]) Get(index int) (V, bool) {
	keys := s.elements
	if index >= 0 && index < len(keys) {
		key := keys[index]
		val, ok := s.uniques[key]
		return val, ok
	}
	var no V
	return no, false
}

func (s *MapValues[K, V]) Filter(filter c.Predicate[V]) c.Pipe[V, []V] {
	return it.NewPipe[V](it.Filter(s.Head(), filter))
}

func (s *MapValues[K, V]) Map(by c.Converter[V, V]) c.Pipe[V, []V] {
	return it.NewPipe[V](it.Map(s.Head(), by))
}

func (s *MapValues[K, V]) Reduce(by op.Binary[V]) V {
	return it.Reduce(s.Head(), by)
}

func (s *MapValues[K, V]) String() string {
	return slice.ToString(s.Collect())
}
