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

func NewSet[T comparable](elements []T) *Set[T] {
	uniques := map[T]struct{}{}
	for _, T := range elements {
		uniques[T] = struct{}{}
	}
	return WrapSet(uniques)
}

func WrapSet[T comparable](uniques map[T]struct{}) *Set[T] {
	return &Set[T]{uniques: uniques}
}

//Set provides uniqueness (does'T insert duplicated values).
type Set[T comparable] struct {
	uniques map[T]struct{}
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ fmt.Stringer = (*Set[int])(nil)
)

func (s *Set[T]) Begin() c.Iterator[T] {
	return s.Head()
}

func (s *Set[T]) Head() *it.Key[T, struct{}] {
	return it.NewKey(s.uniques)
}

func (s *Set[T]) Collect() []T {
	uniques := s.uniques
	out := make([]T, 0, len(uniques))
	for e := range uniques {
		out = append(out, e)
	}
	return out
}

func (s *Set[T]) Len() int {
	return len(s.uniques)
}

func (s *Set[T]) IsEmpty() bool {
	return s.IsEmpty()
}

func (s *Set[T]) For(walker func(T) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	map_.ForEachKey(s.uniques, walker)
}

func (s *Set[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Filter(s.Head(), filter))
}

func (s *Set[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Map(s.Head(), by))
}

func (s *Set[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Head(), by)
}

func (s *Set[T]) Contains(val T) bool {
	_, ok := s.uniques[val]
	return ok
}

func (s *Set[T]) Sort(less func(e1, e2 T) bool) *ordered.Set[T] {
	return ordered.WrapSet(slice.Sort(s.Collect(), less), s.uniques)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.Collect())
}
