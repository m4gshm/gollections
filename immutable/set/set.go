package set

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func Convert[T comparable](elements []T) *Set[T, any] {
	uniques := make(map[T]any, 0)
	for _, v := range elements {
		uniques[v] = struct{}{}
	}
	return Wrap(uniques)
}

func Wrap[k comparable, v any, m map[k]v](uniques m) *Set[k, v] {
	return &Set[k, v]{uniques: uniques}
}

type Set[k comparable, v any] struct {
	uniques map[k]v
}

var _ typ.Set[any, typ.Iterator[any]] = (*Set[any, any])(nil)
var _ fmt.Stringer = (*Set[any, any])(nil)
var _ fmt.GoStringer = (*Set[any, any])(nil)

func (s *Set[k, v]) Begin() typ.Iterator[k] {
	return s.Iter()
}

func (s *Set[k, v]) Iter() *it.Key[k, v] {
	return it.NewKey(s.uniques)
}

func (s *Set[k, v]) Collect() []k {
	uniques := s.uniques
	out := make([]k, 0, len(uniques))
	for e := range uniques {
		out = append(out, e)
	}
	return out
}

func (s *Set[k, v]) For(walker func(k) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *Set[k, v]) ForEach(walker func(k)) error {
	return map_.ForEachKey(s.uniques, walker)
}

func (s *Set[k, v]) Filter(filter typ.Predicate[k]) typ.Pipe[k, []k, typ.Iterator[k]] {
	return it.NewPipe[k](it.Filter(s.Iter(), filter))
}

func (s *Set[k, v]) Map(by typ.Converter[k, k]) typ.Pipe[k, []k, typ.Iterator[k]] {
	return it.NewPipe[k](it.Map(s.Iter(), by))
}

func (s *Set[k, v]) Reduce(by op.Binary[k]) k {
	return it.Reduce(s.Iter(), by)
}

func (s *Set[k, v]) Len() int {
	return len(s.uniques)
}

func (s *Set[k, v]) Contains(val k) bool {
	_, ok := s.uniques[val]
	return ok
}

func (s *Set[k, v]) String() string {
	return s.GoString()
}

func (s *Set[k, v]) GoString() string {
	return slice.ToString(s.Collect())
}
