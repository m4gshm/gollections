package omap

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func WrapVal[k comparable, v any](elements []*k, uniques map[k]v) *MapValues[k, v] {
	return &MapValues[k, v]{elements, uniques}
}

type MapValues[k comparable, v any] struct {
	elements []*k

	uniques map[k]v
}

var _ typ.Collection[any, []any, typ.Iterator[any]] = (*MapValues[any, any])(nil)
var _ fmt.Stringer = (*MapValues[any, any])(nil)

func (s *MapValues[k, v]) Begin() typ.Iterator[v] {
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
	for i, r := range refs {
		key := *r
		val := s.uniques[key]
		elements[i] = val
	}
	return elements
}

func (s *MapValues[k, v]) Track(tracker func(int, v) error) error {
	refs := s.elements
	for i, r := range refs {
		key := *r
		val := s.uniques[key]
		if err := tracker(i, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *MapValues[k, v]) TrackEach(tracker func(int, v)) error {
	return s.Track(func(i int, value v) error { tracker(i, value); return nil })
}

func (s *MapValues[k, v]) For(walker func(v) error) error {
	return s.Track(func(i int, value v) error { return walker(value) })
}

func (s *MapValues[k, v]) ForEach(walker func(v)) error {
	return s.TrackEach(func(_ int, value v) { walker(value) })
}

func (s *MapValues[k, v]) Get(index int) (v, bool) {
	refs := s.elements
	if index >= 0 && index < len(refs) {
		key := *refs[index]
		val, ok := s.uniques[key]
		return val, ok
	}
	var no v
	return no, false
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
