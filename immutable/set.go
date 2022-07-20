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

//NewSet instantiates Set and copies elements to it.
func NewSet[T comparable](elements []T) Set[T] {
	internal := map[T]struct{}{}
	for _, T := range elements {
		internal[T] = struct{}{}
	}
	return WrapSet(internal)
}

//WrapSet creates a set using a map as the internal storage.
func WrapSet[T comparable](elements map[T]struct{}) Set[T] {
	return Set[T]{elements: elements}
}

//Set is the Collection implementation that provides the uniqueness of its elements. Elements must be comparable.
type Set[T comparable] struct {
	elements map[T]struct{}
}

var (
	_ c.Set[int]   = (*Set[int])(nil)
	_ c.Set[int]   = Set[int]{}
	_ fmt.Stringer = (*Set[int])(nil)
	_ fmt.Stringer = Set[int]{}
)

func (s Set[T]) Begin() c.Iterator[T] {
	return s.Head()
}

func (s Set[T]) Head() it.Key[T, struct{}] {
	return it.NewKey(s.elements)
}

func (s Set[T]) First() (it.Key[T, struct{}], T, bool) {
	var (
		iter      = s.Head()
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (s Set[T]) Collect() []T {
	elements := s.elements
	out := make([]T, 0, len(elements))
	for e := range elements {
		out = append(out, e)
	}
	return out
}

func (s Set[T]) Len() int {
	return len(s.elements)
}

func (s Set[T]) IsEmpty() bool {
	return s.IsEmpty()
}

func (s Set[T]) For(walker func(T) error) error {
	return map_.ForKeys(s.elements, walker)
}

func (s Set[T]) ForEach(walker func(T)) {
	map_.ForEachKey(s.elements, walker)
}

func (s Set[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Filter(s.Head(), filter))
}

func (s Set[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	return it.NewPipe[T](it.Map(s.Head(), by))
}

func (s Set[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Head(), by)
}

func (s Set[T]) Contains(val T) bool {
	_, ok := s.elements[val]
	return ok
}

//Sort transforms to the ordered Set.
func (s Set[T]) Sort(less func(e1, e2 T) bool) ordered.Set[T] {
	return ordered.WrapSet(slice.Sort(s.Collect(), less), s.elements)
}

func (s Set[T]) String() string {
	return slice.ToString(s.Collect())
}
