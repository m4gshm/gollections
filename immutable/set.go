package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewSet[T comparable](elements []T) *Set[T] {
	uniques := map[T]struct{}{}
	for _, v := range elements {
		uniques[v] = struct{}{}
	}
	return WrapSet(uniques)
}

func WrapSet[k comparable](uniques map[k]struct{}) *Set[k] {
	return &Set[k]{uniques: uniques}
}

//Set provides uniqueness (does't insert duplicated values).
type Set[k comparable] struct {
	uniques map[k]struct{}
}

var (
	_ c.Set[int] = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s *Set[k]) Begin() c.Iterator[k] {
	return s.Iter()
}

func (s *Set[k]) Iter() *it.Key[k, struct{}] {
	return it.NewKey(s.uniques)
}

func (s *Set[k]) Collect() []k {
	uniques := s.uniques
	out := make([]k, 0, len(uniques))
	for e := range uniques {
		out = append(out, e)
	}
	return out
}

func (s *Set[k]) For(walker func(k) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *Set[k]) ForEach(walker func(k)) {
	map_.ForEachKey(s.uniques, walker)
}

func (s *Set[k]) Filter(filter c.Predicate[k]) c.Pipe[k, []k, c.Iterator[k]] {
	return it.NewPipe[k](it.Filter(s.Iter(), filter))
}

func (s *Set[k]) Map(by c.Converter[k, k]) c.Pipe[k, []k, c.Iterator[k]] {
	return it.NewPipe[k](it.Map(s.Iter(), by))
}

func (s *Set[k]) Reduce(by op.Binary[k]) k {
	return it.Reduce(s.Iter(), by)
}

func (s *Set[k]) Len() int {
	return len(s.uniques)
}

func (s *Set[k]) Contains(val k) bool {
	_, ok := s.uniques[val]
	return ok
}

func (s *Set[k]) String() string {
	return slice.ToString(s.Collect())
}
