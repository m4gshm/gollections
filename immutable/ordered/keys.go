package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

//WrapKeys creates the MapKeys using elements as internal storage.
func WrapKeys[T comparable](elements []T) MapKeys[T] {
	return MapKeys[T]{elements}
}

//MapKeys is the wrapper for Map's keys.
type MapKeys[T comparable] struct {
	elements []T
}

var (
	_ c.Collection[int, []int, c.Iterator[int]] = (*MapKeys[int])(nil)
	_ c.Collection[int, []int, c.Iterator[int]] = MapKeys[int]{}
	_ fmt.Stringer                              = (*MapKeys[int])(nil)
	_ fmt.Stringer                              = MapKeys[int]{}
)

func (s MapKeys[T]) Begin() c.Iterator[T] {
	iter := s.Head()
	return &iter
}

func (s MapKeys[T]) Head() it.Iter[T] {
	return it.NewHead(s.elements)
}

func (s MapKeys[T]) First() (it.Iter[T], T, bool) {
	var (
		iter      = s.Head()
		first, ok = iter.Next()
	)
	return iter, first, ok
}

func (s MapKeys[T]) Len() int {
	return len(s.elements)
}

func (s MapKeys[T]) IsEmpty() bool {
	return s.Len() == 0
}

func (s MapKeys[T]) Collect() []T {
	elements := s.elements
	dest := make([]T, len(elements))
	copy(dest, elements)
	return dest
}

func (s MapKeys[T]) For(walker func(T) error) error {
	return slice.For(s.elements, walker)
}

func (s MapKeys[T]) ForEach(walker func(T)) {
	slice.ForEach(s.elements, walker)
}

func (s MapKeys[T]) Get(index int) (T, bool) {
	return slice.Get(s.elements, index)
}

func (s MapKeys[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T] {
	iter := s.Head()
	return it.NewPipe[T](it.Filter(&iter, filter))
}

func (s MapKeys[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T] {
	iter := s.Head()
	return it.NewPipe[T](it.Map(&iter, by))
}

func (s MapKeys[T]) Reduce(by op.Binary[T]) T {
	iter := s.Head()
	return it.Reduce(&iter, by)
}

func (s MapKeys[T]) String() string {
	return slice.ToString(s.Collect())
}
