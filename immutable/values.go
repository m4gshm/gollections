package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func WrapVal[K comparable, V any](uniques map[K]V) *MapValues[K, V] {
	return &MapValues[K, V]{uniques}
}

type MapValues[K comparable, V any] struct {
	uniques map[K]V
}

var _ c.Collection[any, []any, c.Iterator[any]] = (*MapValues[int, any])(nil)
var _ fmt.Stringer = (*MapValues[int, any])(nil)

func (s *MapValues[K, V]) Begin() c.Iterator[V] {
	return s.Head()
}

func (s *MapValues[K, V]) Head() *it.Val[K, V] {
	return it.NewVal(s.uniques)
}

func (s *MapValues[K, V]) Len() int {
	return len(s.uniques)
}

func (s *MapValues[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *MapValues[K, V]) Collect() []V {
	uniques := s.uniques
	elements := make([]V, 0, len(s.uniques))
	for _, val := range uniques {
		elements = append(elements, val)
	}
	return elements
}

func (s *MapValues[K, V]) For(walker func(V) error) error {
	return map_.ForValues(s.uniques, walker)
}

func (s *MapValues[K, V]) ForEach(walker func(V)) {
	map_.ForEachValue(s.uniques, walker)
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

func (s *MapValues[K, V]) Sort(less func(e1, e2 V) bool) *Vector[V] {
	var dest = s.Collect()
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (s *MapValues[K, V]) String() string {
	return slice.ToString(s.Collect())
}
