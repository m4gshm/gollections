package vector

import (
	"fmt"

	"github.com/m4gshm/container/immutable/vector"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/slice"
	"github.com/m4gshm/container/typ"
)

func Create[T any](capacity int) *Vector[T] {
	return Wrap(make([]T, 0, capacity))
}

func Convert[T any](elements []T) *Vector[T] {
	return Wrap(slice.Copy(elements))
}

func Wrap[T any](elements []T) *Vector[T] {
	return &Vector[T]{Vector: vector.Wrap(elements), elements: &elements}
}

type Vector[t any /*if replaces generic type by 'T' it raises compile-time error 'type parameter bound more than once'*/] struct {
	*vector.Vector[t]
	elements   *[]t
	changeMark int32
}

var _ mutable.Vector[any, typ.Iterator[any]] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[T]) Add(v ...T) (bool, error) {
	return s.AddAll(v)
}

func (s *Vector[T]) AddAll(v []T) (bool, error) {
	markOnStart := s.changeMark
	*s.elements = append(*s.elements, v...)
	return mutable.Commit(markOnStart, &s.changeMark)
}

func (s *Vector[T]) AddOne(v T) (bool, error) {
	markOnStart := s.changeMark
	*s.elements = append(*s.elements, v)
	return mutable.Commit(markOnStart, &s.changeMark)
}
