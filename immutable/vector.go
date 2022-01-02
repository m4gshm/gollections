package immutable

import (
	"fmt"

	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func NewVector[T any](values []T) *Vector[T] {
	elements := make([]T, len(values))
	copy(elements, values)
	return WrapVector(elements)
}

func WrapVector[T any](elements []T) *Vector[T] {
	return &Vector[T]{elements: elements}
}

type Vector[T any] struct {
	elements   []T
	changeMark int32
}

var _ Vec[any] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) Iter() *iter.Iter[T] {
	return iter.New(s.elements)
}

func (s *Vector[T]) Values() []T {
	elements := make([]T, len(s.elements))
	copy(elements, elements)
	return elements
}

func (s *Vector[T]) ForEach(walker func(T)) {
	for _, e := range s.elements {
		walker(e)
	}
}

func (s *Vector[T]) Len() int {
	return len(s.elements)
}

func (s *Vector[T]) Get(index int) (T, bool) {
	elements := s.elements
	l := len(elements)
	if l > 0 && (index >= 0 || index < l) {
		return s.elements[index], true
	}
	var no T
	return no, false
}

func (s *Vector[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T] {
	return iter.NewPipe[T](iter.Filter(s.Iter(), filter))
}

func (s *Vector[T]) Map(by typ.Converter[T, T]) typ.Pipe[T] {
	return iter.NewPipe[T](iter.Map(s.Iter(), by))
}

func (s *Vector[T]) Reduce(by op.Binary[T]) T {
	return iter.Reduce(s.Iter(), by)
}

func (s *Vector[T]) String() string {
	return slice.ToString(s.elements)
}
