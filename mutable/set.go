package mutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/slice"
)

// NewSetCap creates a set with a predefined capacity.
func NewSetCap[T comparable](capacity int) Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

// NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) Set[T] {
	internal := make(map[T]struct{}, len(elements))
	for _, v := range elements {
		internal[v] = struct{}{}
	}
	return WrapSet(internal)
}

// ToSet creates a Set instance with elements obtained by passing an iterator.
func ToSet[T comparable](elements c.Iterator[T]) Set[T] {
	internal := map[T]struct{}{}
	for {
		if e, ok := elements.Next(); !ok {
			break
		} else {
			internal[e] = struct{}{}
		}
	}
	return WrapSet(internal)
}

// WrapSet creates a set using a map as the internal storage.
func WrapSet[K comparable](elements map[K]struct{}) Set[K] {
	return Set[K]{elements: elements}
}

// Set is the Collection implementation that provides element uniqueness. The elements must be comparable.
type Set[K comparable] struct {
	elements map[K]struct{}
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
	_ c.Addable[int]          = Set[int]{}
	_ c.AddableNew[int]       = Set[int]{}
	_ c.AddableAll[int]       = Set[int]{}
	_ c.AddableAllNew[int]    = Set[int]{}
	_ c.Deleteable[int]       = Set[int]{}
	_ c.DeleteableVerify[int] = Set[int]{}
	_ c.Set[int]              = Set[int]{}
	_ fmt.Stringer            = Set[int]{}
)

func (s Set[K]) Begin() c.Iterator[K] {
	return s.Head()
}

func (s Set[K]) BeginEdit() c.DelIterator[K] {
	return s.Head()
}

func (s Set[K]) Head() *SetIter[K] {
	return NewSetIter(s.elements, s.DeleteOne)
}

func (s Set[K]) Collect() []K {
	return map_.Keys(s.elements)
}

func (s Set[K]) Slice() []K {
	return s.Collect()
}

func (s Set[T]) Copy() Set[T] {
	return WrapSet(map_.Clone(s.elements))
}

func (s Set[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Set[K]) Len() int {
	return len(s.elements)
}

func (s Set[K]) Contains(val K) bool {
	_, ok := s.elements[val]
	return ok
}

func (s Set[K]) Add(elements ...K) {
	for _, element := range elements {
		s.AddOne(element)
	}
}

func (s Set[K]) AddOne(element K) {
	s.elements[element] = struct{}{}
}

func (s Set[K]) AddNew(elements ...K) bool {
	ok := false
	for _, element := range elements {
		ok = s.AddOneNew(element) || ok
	}
	return ok
}

func (s Set[K]) AddOneNew(element K) bool {
	ok := !s.Contains(element)
	if ok {
		s.elements[element] = struct{}{}
	}
	return ok
}

func (s Set[T]) AddAll(elements c.Iterable[c.Iterator[T]]) {
	loop.ForEach(elements.Begin().Next, s.AddOne)
}

func (s Set[T]) AddAllNew(elements c.Iterable[c.Iterator[T]]) bool {
	var ok bool
	loop.ForEach(elements.Begin().Next, func(v T) { ok = s.AddOneNew(v) || ok })
	return ok
}

func (s Set[K]) Delete(elements ...K) {
	for _, element := range elements {
		s.DeleteOne(element)
	}
}

func (s Set[K]) DeleteOne(element K) {
	delete(s.elements, element)
}

func (s Set[T]) DeleteActual(elements ...T) bool {
	ok := false
	for i := range elements {
		ok = s.DeleteActualOne(elements[i]) || ok
	}
	return ok
}

func (s Set[K]) DeleteActualOne(element K) bool {
	_, ok := s.elements[element]
	if ok {
		delete(s.elements, element)
	}
	return ok
}

func (s Set[K]) For(walker func(K) error) error {
	return map_.ForKeys(s.elements, walker)
}

func (s Set[K]) ForEach(walker func(K)) {
	map_.ForEachKey(s.elements, walker)
}

func (s Set[K]) Filter(filter func(K) bool) c.Pipe[K, []K] {
	h := s.Head()
	return it.NewPipe[K](it.Filter(h, h.Next, filter))
}

func (s Set[K]) Convert(by func(K) K) c.Pipe[K, []K] {
	h := s.Head()
	return it.NewPipe[K](it.Convert(h, h.Next, by))
}

func (s Set[K]) Reduce(by func(K, K) K) K {
	return loop.Reduce(s.Head().Next, by)
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

func (s Set[K]) String() string {
	return slice.ToString(s.Slice())
}
