package vector

import (
	"fmt"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

func Convert[T any](elements []T) *Vector[T] {
	c := slice.Copy(elements)
	return Wrap(&c)
}

func Wrap[T any](elements *[]T) *Vector[T] {
	return &Vector[T]{elements: elements}
}

type Vector[T any] struct {
	elements *[]T
}

var _ typ.Vector[any, typ.Iterator[any]] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) Iter() *it.Iter[T] {
	return it.New(*s.elements)
}

func (s *Vector[T]) Collect() []T {
	return slice.Copy(*s.elements)
}

func (s *Vector[T]) Track(tracker func(int, T) error) error {
	for i, e := range *s.elements {
		if err := tracker(i, e); err != nil {
			return err
		}
	}
	return nil
}

func (s *Vector[T]) TrackEach(tracker func(int, T)) error {
	return s.Track(func(i int, value T) error { tracker(i, value); return nil })
}

func (s *Vector[T]) For(walker func(T) error) error {
	return s.Track(func(i int, value T) error { return walker(value) })
}

func (s *Vector[T]) ForEach(walker func(T)) error {
	return s.TrackEach(func(_ int, value T) { walker(value) })
}

func (s *Vector[T]) Len() int {
	return len(*s.elements)
}

func (s *Vector[T]) Get(index int) (T, bool) {
	elements := *s.elements
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return elements[index], true
	}
	var no T
	return no, false
}

func (s *Vector[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Vector[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, []T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *Vector[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *Vector[T]) String() string {
	return slice.ToString(*s.elements)
}
