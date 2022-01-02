package mutable

import (
	"fmt"

	"github.com/m4gshm/container/immutable"
)

func NewVector[T any](values []T) *Vector[T] {
	elements := make([]T, len(values))
	copy(elements, values)
	return WrapVector(elements)
}

func WrapVector[T any](elements []T) *Vector[T] {
	return &Vector[T]{Vector: immutable.WrapVector(elements), elements:  &elements}
}

type Vector[T any] struct {
	*immutable.Vector[T]
	elements   *[]T
	changeMark int32
}

var _ Vec[any] = (*Vector[any])(nil)
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


