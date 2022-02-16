package mutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

//NewSet creates a set with a predefined capacity.
func NewSet[T comparable](capacity int) *Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

//ToSet converts an elements slice to the set containing them.
func ToSet[T comparable](elements []T) *Set[T] {
	internal := make(map[T]struct{}, len(elements))
	for _, v := range elements {
		internal[v] = struct{}{}
	}
	return WrapSet(internal)
}

//WrapSet creates a set using a map as the internal storage.
func WrapSet[K comparable](elements map[K]struct{}) *Set[K] {
	return &Set[K]{elements: elements}
}

//Set is the Collection implementation that provides element uniqueness. The elements must be comparable.
type Set[K comparable] struct {
	elements map[K]struct{}
}

var (
	_ c.Addable[int]    = (*Set[int])(nil)
	_ c.Deleteable[int] = (*Set[int])(nil)
	_ c.Set[int]        = (*Set[int])(nil)
	_ fmt.Stringer      = (*Set[int])(nil)
)

func (s *Set[K]) Begin() c.Iterator[K] {
	return s.Head()
}

func (s *Set[K]) BeginEdit() c.DelIterator[K] {
	return s.Head()
}

func (s *Set[K]) Head() *SetIter[K] {
	return NewSetIter(s.elements, s.DeleteOne)
}

func (s *Set[K]) Collect() []K {
	return map_.Keys(s.elements)
}

func (s *Set[T]) Copy() *Set[T] {
	return WrapSet(map_.Copy(s.elements))
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[K]) Len() int {
	return len(s.elements)
}

func (s *Set[K]) Contains(val K) bool {
	_, ok := s.elements[val]
	return ok
}

func (s *Set[K]) Add(elements ...K) bool {
	return s.AddAll(elements)
}

func (s *Set[K]) AddAll(elements []K) bool {
	added := false
	for _, element := range elements {
		if _, ok := s.elements[element]; !ok {
			s.elements[element] = struct{}{}
			added = true
		}
	}
	return added
}

func (s *Set[K]) AddOne(element K) bool {
	if _, ok := s.elements[element]; ok {
		return false
	}
	s.elements[element] = struct{}{}
	return true
}

func (s *Set[K]) Delete(elements ...K) bool {
	for _, element := range elements {
		if _, ok := s.elements[element]; !ok {
			return false
		}
		delete(s.elements, element)
	}
	return true
}

func (s *Set[K]) DeleteOne(element K) bool {
	if _, ok := s.elements[element]; !ok {
		return false
	}
	delete(s.elements, element)
	return true
}

func (s *Set[K]) For(walker func(K) error) error {
	return map_.ForKeys(s.elements, walker)
}

func (s *Set[K]) ForEach(walker func(K)) {
	map_.ForEachKey(s.elements, walker)
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

//Sort transforms to the ordered Set contains sorted elements.
func (s *Set[T]) Sort(less func(e1, e2 T) bool) *ordered.Set[T] {
	return ordered.ToSet(slice.Sort(s.Collect(), less))
}

func (s *Set[K]) String() string {
	return slice.ToString(s.Collect())
}
