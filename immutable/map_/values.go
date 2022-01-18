package map_

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func WrapVal[k comparable, v any](uniques map[k]v) *MapValues[k, v] {
	return &MapValues[k, v]{uniques}
}

type MapValues[k comparable, v any] struct {
	uniques map[k]v
}

var _ typ.Collection[any, []any, typ.Iterator[any]] = (*MapValues[any, any])(nil)
var _ fmt.Stringer = (*MapValues[any, any])(nil)

func (s *MapValues[k, v]) Begin() typ.Iterator[v] {
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
	for _, val := range s.uniques {
		if err := walker(val); err != nil {
			return err
		}
	}
	return nil
}

func (s *MapValues[k, v]) ForEach(walker func(v)) error {
	return s.For(func(value v) error { walker(value); return nil })
}

func (s *MapValues[k, v]) Filter(filter typ.Predicate[v]) typ.Pipe[v, []v, typ.Iterator[v]] {
	return it.NewPipe[v](it.Filter(s.Iter(), filter))
}

func (s *MapValues[k, v]) Map(by typ.Converter[v, v]) typ.Pipe[v, []v, typ.Iterator[v]] {
	return it.NewPipe[v](it.Map(s.Iter(), by))
}

func (s *MapValues[k, v]) Reduce(by op.Binary[v]) v {
	return it.Reduce(s.Iter(), by)
}

func (s *MapValues[k, v]) String() string {
	return slice.ToString(s.Collect())
}
