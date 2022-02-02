package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func NewVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Copy(elements))
}

func WrapVector[T any](elements []T) *Vector[T] {
	return &Vector[T]{elements: elements}
}

//Vector stores ordered elements, provides index access.
type Vector[T any] struct {
	elements []T
}

var (
	_ c.Vector[any] = (*Vector[any])(nil)
	_ fmt.Stringer    = (*Vector[any])(nil)
)

func (s *Vector[T]) Begin() c.Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) Iter() *it.PIter[T] {
	return it.NewP(&s.elements)
}

func (s *Vector[T]) Collect() []T {
	return slice.Copy(s.elements)
}

func (s *Vector[T]) Track(tracker func(int, T) error) error {
	return slice.Track(s.elements, tracker)
}

func (s *Vector[T]) TrackEach(tracker func(int, T)) {
	slice.TrackEach(s.elements, tracker)
}

func (s *Vector[T]) For(walker func(T) error) error {
	return slice.For(s.elements, walker)
}

func (s *Vector[T]) ForEach(walker func(T)) {
	slice.ForEach(s.elements, walker)
}

func (s *Vector[T]) Len() int {
	return it.GetLen(&s.elements)
}

func (s *Vector[T]) Get(index int) (T, bool) {
	return slice.Get(s.elements, index)

}

func (s *Vector[T]) Filter(filter c.Predicate[T]) c.Pipe[T, []T, c.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Vector[T]) Map(by c.Converter[T, T]) c.Pipe[T, []T, c.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *Vector[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *Vector[T]) String() string {
	return slice.ToString(s.elements)
}
