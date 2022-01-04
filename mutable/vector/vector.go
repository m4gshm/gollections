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
	return &Vector[T]{Vector: vector.Wrap(elements), elements:  &elements}
}

type Vector[T any] struct {
	*vector.Vector[T]
	elements   *[]T
	changeMark int32
}

var _ mutable.Vector[any, typ.Iterator[any]] = (*Vector[any])(nil)
var _ fmt.Stringer = (*Vector[any])(nil)

func (s *Vector[T]) Add(v T) bool {
	markOnStart := s.changeMark
	*s.elements = append(*s.elements, v)
	markOnFinish := s.changeMark
	if markOnFinish != markOnStart {
		panic("concurrent Slice read and write")
	}
	return true
}


