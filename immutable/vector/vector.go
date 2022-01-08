package vector

import (
	"fmt"

	"github.com/m4gshm/container/immutable"
	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/op"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
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

var _ immutable.Vector[any, typ.Iterator[any]] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[T]) Begin() typ.Iterator[T] {
	return s.Iter()
}

func (s *Vector[T]) Iter() *it.Iter[T] {
	return it.New(*s.elements)
}

func (s *Vector[T]) Elements() []T {
	return slice.Copy(*s.elements)
}

func (s *Vector[T]) ForEach(walker func(T)) {
	for _, e := range *s.elements {
		walker(e)
	}
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

func (s *Vector[T]) Filter(filter typ.Predicate[T]) typ.Pipe[T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Filter(s.Iter(), filter))
}

func (s *Vector[T]) Map(by typ.Converter[T, T]) typ.Pipe[T, typ.Iterator[T]] {
	return it.NewPipe[T](it.Map(s.Iter(), by))
}

func (s *Vector[T]) Reduce(by op.Binary[T]) T {
	return it.Reduce(s.Iter(), by)
}

func (s *Vector[T]) String() string {
	return slice.ToString(*s.elements)
}
