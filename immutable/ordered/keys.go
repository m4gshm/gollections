package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// WrapKeys instantiates MapKeys using elements as internal storage
func WrapKeys[T comparable](elements []T) MapKeys[T] {
	return MapKeys[T]{elements}
}

// MapKeys is the wrapper for Map's keys
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
	h := s.Head()
	return &h
}

func (s MapKeys[T]) Head() iter.ArrayIter[T] {
	return iter.NewHead(s.elements)
}

func (s MapKeys[T]) First() (iter.ArrayIter[T], T, bool) {
	var (
		iterator  = s.Head()
		first, ok = iterator.Next()
	)
	return iterator, first, ok
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

func (s MapKeys[T]) Filter(filter func(T) bool) c.Pipe[T, []T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Filter(h, h.Next, filter))
}

func (s MapKeys[T]) Convert(by func(T) T) c.Pipe[T, []T] {
	h := s.Head()
	return iter.NewPipe[T](iter.Convert(h, h.Next, by))
}

func (s MapKeys[T]) Reduce(by func(T, T) T) T {
	h := s.Head()
	return loop.Reduce(h.Next, by)
}

func (s MapKeys[T]) String() string {
	return slice.ToString(s.Collect())
}
