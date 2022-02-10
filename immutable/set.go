package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewSet[t comparable](elements []t) *Set[t] {
	uniques := map[t]struct{}{}
	for _, t := range elements {
		uniques[t] = struct{}{}
	}
	return WrapSet(uniques)
}

func WrapSet[t comparable](uniques map[t]struct{}) *Set[t] {
	return &Set[t]{uniques: uniques}
}

//Set provides uniqueness (does't insert duplicated values).
type Set[t comparable] struct {
	uniques map[t]struct{}
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s *Set[t]) Begin() c.Iterator[t] {
	return s.Iter()
}

func (s *Set[t]) Iter() *it.Key[t, struct{}] {
	return it.NewKey(s.uniques)
}

func (s *Set[t]) Collect() []t {
	uniques := s.uniques
	out := make([]t, 0, len(uniques))
	for e := range uniques {
		out = append(out, e)
	}
	return out
}

func (s *Set[t]) Len() int {
	return len(s.uniques)
}

func (s *Set[t]) For(walker func(t) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *Set[t]) ForEach(walker func(t)) {
	map_.ForEachKey(s.uniques, walker)
}

func (s *Set[t]) Filter(filter c.Predicate[t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Filter(s.Iter(), filter))
}

func (s *Set[t]) Map(by c.Converter[t, t]) c.Pipe[t, []t, c.Iterator[t]] {
	return it.NewPipe[t](it.Map(s.Iter(), by))
}

func (s *Set[t]) Reduce(by op.Binary[t]) t {
	return it.Reduce(s.Iter(), by)
}

func (s *Set[t]) Contains(val t) bool {
	_, ok := s.uniques[val]
	return ok
}

func (s *Set[t]) Sort(less func(e1, e2 t) bool) *ordered.Set[t] {
	return ordered.WrapSet(slice.Sort(s.Collect(), less), s.uniques)
}

func (s *Set[t]) String() string {
	return slice.ToString(s.Collect())
}
