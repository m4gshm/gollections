package dict

import (
	"fmt"

	"github.com/m4gshm/container/immutable"
	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func Wrap[k comparable, v any](elements []*k, uniques map[k]v) *Vector[k, v] {
	return &Vector[k, v]{elements, uniques}
}

type Vector[k comparable, v any] struct {
	elements []*k
	uniques  map[k]v
}

var _ immutable.Vector[any, typ.Iterator[any]] = (*Vector[any, any])(nil)
var _ fmt.Stringer = (*Vector[any, any])(nil)

func (s *Vector[k, v]) Begin() typ.Iterator[v] {
	return s.Iter()
}

func (s *Vector[k, v]) Iter() *Iterator[k, v] {
	return NewIterator(s.elements, s.uniques)
}

func (s *Vector[k, v]) Len() int {
	return len(s.elements)
}

func (s *Vector[k, v]) Elements() []v {
	refs := s.elements
	elements := make([]v, len(refs))
	for i, r := range refs {
		key := *r
		val := s.uniques[key]
		elements[i] = val
	}
	return elements
}

func (s *Vector[k, v]) ForEach(walker func(v)) {
	refs := s.elements
	for _, r := range refs {
		key := *r
		val := s.uniques[key]
		walker(val)
	}
}

func (s *Vector[k, v]) Get(index int) (v, bool) {
	refs := s.elements
	if index >= 0 && index < len(refs) {
		key := *refs[index]
		val, ok := s.uniques[key]
		return val, ok
	}
	var no v
	return no, false
}

func (s *Vector[k, v]) Filter(filter typ.Predicate[v]) typ.Pipe[v, typ.Iterator[v]] {
	return it.NewPipe[v](it.Filter(s.Iter(), filter))
}

func (s *Vector[k, v]) Map(by typ.Converter[v, v]) typ.Pipe[v, typ.Iterator[v]] {
	return it.NewPipe[v](it.Map(s.Iter(), by))
}

func (s *Vector[k, v]) Reduce(by op.Binary[v]) v {
	return it.Reduce(s.Iter(), by)
}

func (s *Vector[k, v]) String() string {
	return slice.ToString(s.Elements())
}
