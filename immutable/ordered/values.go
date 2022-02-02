package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func WrapVal[k comparable, v any](elements []k, uniques map[k]v) *MapValues[k, v] {
	return &MapValues[k, v]{elements, uniques}
}

type MapValues[k comparable, v any] struct {
	elements []k

	uniques map[k]v
}

var _ c.Collection[any, []any, c.Iterator[any]] = (*MapValues[int, any])(nil)
var _ fmt.Stringer = (*MapValues[int, any])(nil)

func (s *MapValues[k, v]) Begin() c.Iterator[v] {
	return s.Iter()
}

func (s *MapValues[k, v]) Iter() *ValIter[k, v] {
	return NewValIter(s.elements, s.uniques)
}

func (s *MapValues[k, v]) Len() int {
	return len(s.uniques)
}

func (s *MapValues[k, v]) Collect() []v {
	refs := s.elements
	elements := make([]v, len(refs))
	for i, key := range refs {
		val := s.uniques[key]
		elements[i] = val
	}
	return elements
}

func (s *MapValues[k, v]) For(walker func(v) error) error {
	return map_.ForOrderedValues(s.elements, s.uniques, walker)
}

func (s *MapValues[k, v]) ForEach(walker func(v)) {
	map_.ForEachOrderedValues(s.elements, s.uniques, walker)
}

func (s *MapValues[k, v]) Get(index int) (v, bool) {
	keys := s.elements
	if index >= 0 && index < len(keys) {
		key := keys[index]
		val, ok := s.uniques[key]
		return val, ok
	}
	var no v
	return no, false
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

func (s *MapValues[k, v]) String() string {
	return slice.ToString(s.Collect())
}
