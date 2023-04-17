package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) *Set[T] {
	var (
		l       = len(elements)
		uniques = make(map[T]int, l)
		order   = make([]T, 0, l)
	)
	pos := 0
	for _, e := range elements {
		if _, ok := uniques[e]; !ok {
			order = append(order, e)
			uniques[e] = pos
			pos++
		}
	}
	return WrapSet(order, uniques)
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) *Set[T] {
	var (
		uniques = map[T]int{}
		order   []T
		pos     = 0
	)

	for {
		if e, ok := elements.Next(); !ok {
			break
		} else {
			order, pos = add(e, uniques, order, pos)
		}
	}
	return WrapSet(order, uniques)
}

// NewSetCap creates a set with a predefined capacity.
func NewSetCap[T comparable](capacity int) *Set[T] {
	return WrapSet(make([]T, 0, capacity), make(map[T]int, capacity))
}

// WrapSet creates a set using a map and an order slice as the internal storage.
func WrapSet[T comparable](elements []T, uniques map[T]int) *Set[T] {
	return &Set[T]{elements: elements, uniques: uniques}
}

// Set is the Collection implementation that provides element uniqueness and access order. Elements must be comparable.
type Set[T comparable] struct {
	elements []T
	uniques  map[T]int
}

var (
	_ c.Addable[int]          = (*Set[int])(nil)
	_ c.AddableNew[int]       = (*Set[int])(nil)
	_ c.Deleteable[int]       = (*Set[int])(nil)
	_ c.DeleteableVerify[int] = (*Set[int])(nil)
	_ c.Set[int]              = (*Set[int])(nil)
	_ fmt.Stringer            = (*Set[int])(nil)
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
	return slice.Clone(s.elements)
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

func (s *Set[T]) AddNew(elements ...T) bool {
	ok := false
	for i := range elements {
		ok = s.AddNewOne(elements[i]) || ok
	}
	return ok
}

func (s *Set[T]) AddNewOne(v T) bool {
	u := s.uniques
	if _, ok := u[v]; !ok {
		e := s.elements
		u[v] = len(e)
		s.elements = append(e, v)
		return true
	}
	return false
}

func (s *Set[T]) Add(elements ...T) {
	s.AddNew(elements...)
}

func (s *Set[T]) AddOne(v T) {
	s.AddNewOne(v)
}

func (s *Set[T]) Delete(elements ...T) {
	s.DeleteActual(elements...)
}

func (s *Set[T]) DeleteOne(v T) {
	s.DeleteActualOne(v)
}

func (s *Set[T]) DeleteActual(elements ...T) bool {
	ok := false
	for i := range elements {
		ok = s.DeleteActualOne(elements[i]) || ok
	}
	return ok
}

func (s *Set[T]) DeleteActualOne(v T) bool {
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

func (s *Set[T]) Filter(filter func(T) bool) c.Pipe[T, []T] {
	h := s.Head()
	return it.NewPipe[T](it.Filter(h, h.Next, filter))
}

func (s *Set[T]) Convert(by func(T) T) c.Pipe[T, []T] {
	h := s.Head()
	return it.NewPipe[T](it.Convert(h, h.Next, by))
}

func (s *Set[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(s.Head().Next, by)
}

// Sort transforms to the ordered Set.
func (s *Set[T]) Sort(less slice.Less[T]) *Set[T] {
	return s.sortBy(sort.Slice, less)
}

func (s *Set[T]) StableSort(less slice.Less[T]) *Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s *Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *Set[T] {
	slice.Sort(s.elements, sorter, less)
	return s
}

func (s *Set[T]) String() string {
	return slice.ToString(s.elements)
}

func add[T comparable](e T, uniques map[T]int, order []T, pos int) ([]T, int) {
	if _, ok := uniques[e]; !ok {
		order = append(order, e)
		uniques[e] = pos
		pos++
	}
	return order, pos
}
