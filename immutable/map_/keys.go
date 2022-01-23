package map_

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func WrapKeys[k comparable, v any](uniques map[k]v) *MapKeys[k, v] {
	return &MapKeys[k, v]{uniques}
}

type MapKeys[k comparable, v any] struct {
	uniques map[k]v
}

var _ typ.Collection[any, []any, typ.Iterator[any]] = (*MapKeys[any, any])(nil)
var _ fmt.Stringer = (*MapValues[any, any])(nil)

func (s *MapKeys[k, v]) Begin() typ.Iterator[k] {
	return s.Iter()
}

func (s *MapKeys[k, v]) Iter() *it.Key[k, v] {
	return it.NewKey(s.uniques)
}

func (s *MapKeys[k, v]) Len() int {
	return len(s.uniques)
}

func (s *MapKeys[k, v]) Collect() []k {
	uniques := s.uniques
	elements := make([]k, 0, len(s.uniques))
	for key, _ := range uniques {
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

func (s *MapKeys[k, v]) Filter(filter typ.Predicate[k]) typ.Pipe[k, []k, typ.Iterator[k]] {
	return it.NewPipe[k](it.Filter(s.Iter(), filter))
}

func (s *MapKeys[k, v]) Map(by typ.Converter[k, k]) typ.Pipe[k, []k, typ.Iterator[k]] {
	return it.NewPipe[k](it.Map(s.Iter(), by))
}

func (s *MapKeys[k, v]) Reduce(by op.Binary[k]) k {
	return it.Reduce(s.Iter(), by)
}

func (s *MapKeys[k, v]) String() string {
	return slice.ToString(s.Collect())
}
