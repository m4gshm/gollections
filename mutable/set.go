package mutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/slice"
)

// NewSetCap creates a set with a predefined capacity.
func NewSetCap[T comparable](capacity int) *Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) *Set[T] {
	internal := make(map[T]struct{}, len(elements))
	for _, v := range elements {
		internal[v] = struct{}{}
	}
	return WrapSet(internal)
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) *Set[T] {
	internal := map[T]struct{}{}
	if elements != nil {
		for e, ok := elements.Next(); ok; e, ok = elements.Next() {
			internal[e] = struct{}{}
		}
	}
	return WrapSet(internal)
}

// WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) *Set[T] {
	return &Set[T]{elements: elements}
}

// Set is the Collection implementation that provides element uniqueness. The elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ c.Addable[int]          = (*Set[int])(nil)
	_ c.AddableNew[int]       = (*Set[int])(nil)
	_ c.AddableAll[int]       = (*Set[int])(nil)
	_ c.AddableAllNew[int]    = (*Set[int])(nil)
	_ c.Deleteable[int]       = (*Set[int])(nil)
	_ c.DeleteableVerify[int] = (*Set[int])(nil)
	_ c.Set[int]              = (*Set[int])(nil)
	_ fmt.Stringer            = (*Set[int])(nil)
	// _ c.Addable[int]          = Set[int]{}
	// _ c.AddableNew[int]       = Set[int]{}
	// _ c.AddableAll[int]       = Set[int]{}
	// _ c.AddableAllNew[int]    = Set[int]{}
	// _ c.Deleteable[int]       = Set[int]{}
	// _ c.DeleteableVerify[int] = Set[int]{}
	// _ c.Set[int]              = Set[int]{}
	// _ fmt.Stringer            = Set[int]{}
)

func (s *Set[T]) Begin() c.Iterator[T] {
	h := s.Head()
	return &h
}

func (s *Set[T]) BeginEdit() c.DelIterator[T] {
	h := s.Head()
	return &h
}

func (s *Set[T]) Head() SetIter[T] {
	if s == nil || s.elements == nil {
		return NewSetIter[T](nil, nil)
	}
	return NewSetIter(s.elements, s.DeleteOne)
}

func (s *Set[T]) First() (SetIter[T], T, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
}

func (s *Set[T]) Slice() []T {
	if s == nil || s.elements == nil {
		return nil
	}
	return map_.Keys(s.elements)
}

func (s *Set[T]) Copy() *Set[T] {
	if s == nil || s.elements == nil {
		return nil
	}
	return WrapSet(map_.Clone(s.elements))
}

func (s *Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set[T]) Len() int {
	if s == nil || s.elements == nil {
		return 0
	}
	return len(s.elements)
}

func (s *Set[T]) Contains(val T) bool {
	if s == nil || s.elements == nil {
		return false
	}
	_, ok := s.elements[val]
	return ok
}

func (s *Set[T]) Add(elements ...T) {
	if s == nil {
		return
	}
	for _, element := range elements {
		s.AddOne(element)
	}
}

func (s *Set[T]) AddOne(element T) {
	if s == nil {
		return
	} else if elements := s.elements; elements == nil {
		s.elements = map[T]struct{}{element: {}}
	} else {
		elements[element] = struct{}{}
	}
}

func (s *Set[T]) AddNew(elements ...T) bool {
	if s == nil {
		return false
	}
	ok := false
	for _, element := range elements {
		ok = s.AddOneNew(element) || ok
	}
	return ok
}

func (s *Set[T]) AddOneNew(element T) bool {
	if s == nil {
		return false
	}
	ok := !s.Contains(element)
	if ok {
		s.elements[element] = struct{}{}
	}
	return ok
}

func (s *Set[T]) AddAll(elements c.Iterable[T]) {
	if s == nil {
		return
	}
	loop.ForEach(elements.Begin().Next, s.AddOne)
}

func (s *Set[T]) AddAllNew(elements c.Iterable[T]) bool {
	if s == nil {
		return false
	}
	var ok bool
	loop.ForEach(elements.Begin().Next, func(v T) { ok = s.AddOneNew(v) || ok })
	return ok
}

func (s *Set[T]) Delete(elements ...T) {
	if s == nil {
		return
	}
	for _, element := range elements {
		s.DeleteOne(element)
	}
}

func (s *Set[T]) DeleteOne(element T) {
	if s == nil {
		return
	}
	delete(s.elements, element)
}

func (s *Set[T]) DeleteActual(elements ...T) bool {
	if s == nil {
		return false
	}
	ok := false
	for i := range elements {
		ok = s.DeleteActualOne(elements[i]) || ok
	}
	return ok
}

func (s *Set[T]) DeleteActualOne(element T) bool {
	if s == nil || s.elements == nil {
		return false
	}
	_, ok := s.elements[element]
	if ok {
		delete(s.elements, element)
	}
	return ok
}

func (s *Set[T]) For(walker func(T) error) error {
	if s == nil {
		return nil
	}
	return map_.ForKeys(s.elements, walker)
}

func (s *Set[T]) ForEach(walker func(T)) {
	if s != nil {
		map_.ForEachKey(s.elements, walker)
	}
}

func (s *Set[T]) Filter(filter func(T) bool) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Filter(h, h.Next, filter))
}

func (s *Set[T]) Convert(by func(T) T) c.Pipe[T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Convert(h, h.Next, by))
}

func (s *Set[T]) Reduce(by func(T, T) T) T {
	h := s.Head()
	return loop.Reduce(h.Next, by)
}

// Sort transforms to the ordered Set contains sorted elements.
func (s *Set[T]) Sort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.Slice, less)
}

func (s *Set[T]) StableSort(less slice.Less[T]) *ordered.Set[T] {
	return s.sortBy(sort.SliceStable, less)
}

func (s *Set[T]) sortBy(sorter slice.Sorter, less slice.Less[T]) *ordered.Set[T] {
	cl := slice.Clone(s.Slice())
	slice.Sort(cl, sorter, less)
	return ordered.NewSet(cl)
}

func (s *Set[T]) String() string {
	return slice.ToString(s.Slice())
}
