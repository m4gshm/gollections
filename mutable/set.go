package mutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewSet[T comparable](capacity int) *Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

func ToSet[T comparable](elements []T) *Set[T] {
	uniques := make(map[T]struct{}, len(elements))
	for _, v := range elements {
		uniques[v] = struct{}{}
	}
	return WrapSet(uniques)
}

func WrapSet[K comparable](uniques map[K]struct{}) *Set[K] {
	return &Set[K]{uniques: uniques}
}

//Set provides uniqueness (does't insert duplicated values).
type Set[K comparable] struct {
	uniques map[K]struct{}
}

var (
	_ Addable[int]    = (*Set[int])(nil)
	_ Deleteable[int] = (*Set[int])(nil)
	_ c.Set[int]      = (*Set[int])(nil)
	_ fmt.Stringer    = (*Set[int])(nil)
)

func (s *Set[K]) Begin() c.Iterator[K] {
	return s.Head()
}

func (s *Set[K]) BeginEdit() Iterator[K] {
	return s.Head()
}

func (s *Set[K]) Head() *SetIter[K] {
	return NewSetIter(s.uniques, s.DeleteOne)
}

func (s *Set[K]) Collect() []K {
	uniques := s.uniques
	out := make([]K, 0, len(uniques))
	for e := range uniques {
		out = append(out, e)
	}
	return out
}

func (s *Set[T]) Copy() *Set[T] {
	return WrapSet(map_.Copy(s.uniques))
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[K]) Len() int {
	return len(s.uniques)
}

func (s *Set[K]) Contains(val K) bool {
	_, ok := s.uniques[val]
	return ok
}

func (s *Set[K]) Add(elements ...K) bool {
	return s.AddAll(elements)
}

func (s *Set[K]) AddAll(elements []K) bool {
	uniques := s.uniques
	added := false
	for _, element := range elements {
		if _, ok := uniques[element]; !ok {
			uniques[element] = struct{}{}
			added = true
		}
	}
	return added
}

func (s *Set[K]) AddOne(element K) bool {
	uniques := s.uniques
	if _, ok := uniques[element]; ok {
		return false
	}
	uniques[element] = struct{}{}
	return true
}

func (s *Set[K]) Delete(elements ...K) bool {
	uniques := s.uniques
	for _, element := range elements {
		if _, ok := uniques[element]; !ok {
			return false
		}
		delete(uniques, element)
	}
	return true
}

func (s *Set[K]) DeleteOne(element K) bool {
	uniques := s.uniques
	if _, ok := uniques[element]; !ok {
		return false
	}
	delete(uniques, element)
	return true
}

func (s *Set[K]) For(walker func(K) error) error {
	return map_.ForKeys(s.uniques, walker)
}

func (s *Set[K]) ForEach(walker func(K)) {
	map_.ForEachKey(s.uniques, walker)
}

func (s *Set[K]) Filter(filter c.Predicate[K]) c.Pipe[K, []K] {
	return it.NewPipe[K](it.Filter(s.Head(), filter))
}

func (s *Set[K]) Map(by c.Converter[K, K]) c.Pipe[K, []K] {
	return it.NewPipe[K](it.Map(s.Head(), by))
}

func (s *Set[K]) Reduce(by op.Binary[K]) K {
	return it.Reduce(s.Head(), by)
}

func (s *Set[K]) String() string {
	return slice.ToString(s.Collect())
}
