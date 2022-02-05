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

func WrapVal[k comparable, v any](uniques map[k]v) *MapValues[k, v] {
	return &MapValues[k, v]{uniques}
}

type MapValues[k comparable, v any] struct {
	uniques map[k]v
}

var _ c.Collection[any, []any, c.Iterator[any]] = (*MapValues[int, any])(nil)
var _ fmt.Stringer = (*MapValues[int, any])(nil)

func (s *MapValues[k, v]) Begin() c.Iterator[v] {
	return s.Iter()
}

func (s *MapValues[k, v]) Iter() *it.Val[k, v] {
	return it.NewVal(s.uniques)
}

func (s *MapValues[k, v]) Len() int {
	return len(s.uniques)
}

func (s *MapValues[k, v]) Collect() []v {
	uniques := s.uniques
	elements := make([]v, 0, len(s.uniques))
	for _, val := range uniques {
		elements = append(elements, val)
	}
	return elements
}

func (s *MapValues[k, v]) For(walker func(v) error) error {
	return map_.ForValues(s.uniques, walker)
}

func (s *MapValues[k, v]) ForEach(walker func(v)) {
	map_.ForEachValue(s.uniques, walker)
}

func (s *MapValues[k, v]) Filter(filter c.Predicate[v]) c.Pipe[v, []v, c.Iterator[v]] {
	return it.NewPipe[v](it.Filter(s.Iter(), filter))
}

func (s *MapValues[k, v]) Map(by c.Converter[v, v]) c.Pipe[v, []v, c.Iterator[v]] {
	return it.NewPipe[v](it.Map(s.Iter(), by))
}

func (s *MapValues[k, v]) Reduce(by op.Binary[v]) v {
	return it.Reduce(s.Iter(), by)
}

func (s *MapValues[k, v]) Sort(less func(e1, e2 v) bool) *Vector[v] {
	var dest = s.Collect()
	sort.Slice(dest, func(i, j int) bool { return less(dest[i], dest[j]) })
	return WrapVector(dest)
}

func (s *MapValues[k, v]) String() string {
	return slice.ToString(s.Collect())
}
