package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

//WrapKeys is non-copy constructor
func WrapKeys[k comparable, v any](uniques map[k]v) *MapKeys[k, v] {
	return &MapKeys[k, v]{uniques}
}

//MapKeys is the container reveal keys of a map and hides values.
type MapKeys[k comparable, v any] struct {
	uniques map[k]v
}

var _ c.Collection[int, []int, c.Iterator[int]] = (*MapKeys[int, any])(nil)
var _ fmt.Stringer = (*MapValues[int, any])(nil)

func (s *MapKeys[k, v]) Begin() c.Iterator[k] {
	return s.Head()
}

func (s *MapKeys[k, v]) Head() *it.Key[k, v] {
	return it.NewKey(s.uniques)
}

func (s *MapKeys[k, v]) Len() int {
	return len(s.uniques)
}

func (s *MapKeys[k, v]) Collect() []k {
	uniques := s.uniques
	elements := make([]k, 0, len(s.uniques))
	for key := range uniques {
		elements = append(elements, key)
	}
	return elements
}

func (s *MapKeys[k, v]) For(walker func(k) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *MapKeys[k, v]) ForEach(walker func(k)) {
	map_.ForEachKey(s.uniques, walker)
}

func (s *MapKeys[k, v]) Filter(filter c.Predicate[k]) c.Pipe[k, []k, c.Iterator[k]] {
	return it.NewPipe[k](it.Filter(s.Head(), filter))
}

func (s *MapKeys[k, v]) Map(by c.Converter[k, k]) c.Pipe[k, []k, c.Iterator[k]] {
	return it.NewPipe[k](it.Map(s.Head(), by))
}

func (s *MapKeys[k, v]) Reduce(by op.Binary[k]) k {
	return it.Reduce(s.Head(), by)
}

func (s *MapKeys[k, v]) String() string {
	return slice.ToString(s.Collect())
}
