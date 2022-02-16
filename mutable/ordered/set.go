package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

//ToSet converts an elements slice to the set containing them.
func ToSet[T comparable](elements []T) *Set[T] {
	var (
		l       = len(elements)
		uniques = make(map[T]int, l)
		order   = make([]T, 0, l)
	)
	pos := 0
	for _, v := range elements {
		if _, ok := uniques[v]; !ok {
			order = append(order, v)
			uniques[v] = pos
			pos++
		}
	}
	return WrapSet(order, uniques)
}

//NewSet creates a set with a predefined capacity.
func NewSet[T comparable](capacity int) *Set[T] {
	return WrapSet(make([]T, 0, capacity), make(map[T]int, capacity))
}

//WrapSet creates a set using a map and an order slice as the internal storage.
func WrapSet[T comparable](elements []T, uniques map[T]int) *Set[T] {
	return &Set[T]{elements: elements, uniques: uniques}
}

//Set is the Collection implementation that provides element uniqueness and access order. Elements must be comparable.
type Set[T comparable] struct {
	elements []T
	uniques  map[T]int
}

var (
	_ c.Addable[int]    = (*Set[int])(nil)
	_ c.Deleteable[int] = (*Set[int])(nil)
	_ c.Set[int]        = (*Set[int])(nil)
	_ fmt.Stringer      = (*Set[int])(nil)
)

func (s *Set[T]) Begin() c.Iterator[T] {
	return s.Head()
}

func (s *Set[T]) BeginEdit() c.DelIterator[T] {
	return s.Head()
}

func (s *Set[T]) Head() *SetIter[T] {
	return NewSetIter(&s.elements, s.DeleteOne)
}

func (s *Set[T]) Collect() []T {
	return slice.Copy(s.elements)
}

func (s *Set[T]) For(walker func(T) error) error {
	return slice.For(s.elements, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	slice.ForEach(s.elements, walker)
}

func (s *Set[T]) Len() int {
	return len(s.elements)
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Contains(v T) bool {
	_, ok := s.uniques[v]
	return ok
}

func (s *Set[T]) Add(elements ...T) bool {
	return s.AddAll(elements)
}

func (s *Set[T]) AddAll(elements []T) bool {
	u := s.uniques
	result := false
	for i := range elements {
		v := elements[i]
		if _, ok := u[v]; !ok {
			e := s.elements
			u[v] = len(e)
			s.elements = append(e, v)
			result = true
		}
	}
	return result
}

func (s *Set[T]) AddOne(v T) bool {
	u := s.uniques
	if _, ok := u[v]; !ok {
		e := s.elements
		u[v] = len(e)
		s.elements = append(e, v)
		return true
	}
	return false
}

func (s *Set[T]) Delete(elements ...T) bool {
	return s.DeleteAll(elements)
}

func (s *Set[T]) DeleteAll(elements []T) bool {
	u := s.uniques
	result := false
	for i := range elements {
		v := elements[i]
		if pos, ok := u[v]; ok {
			delete(u, v)
			//todo: need optimize
			e := s.elements
			ne := slice.Delete(pos, e)
			for i := pos; i < len(ne); i++ {
				u[ne[i]]--
			}
			s.elements = ne
			result = true
		}
	}
	return result
}

func (s *Set[T]) DeleteOne(v T) bool {
	u := s.uniques
	if pos, ok := u[v]; ok {
		delete(u, v)
		//todo: need optimize
		e := s.elements
		ne := slice.Delete(pos, e)
		for i := pos; i < len(ne); i++ {
			u[ne[i]]--
		}
		s.elements = ne
		return true
	}
	return false
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

func (s *Set[t]) Sort(less func(e1, e2 t) bool) *Set[t] {
	s.elements = slice.Sort(s.elements, less)
	return s
}

func (s *Set[T]) String() string {
	return slice.ToString(s.elements)
}
