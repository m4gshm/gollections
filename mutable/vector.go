package mutable

import (
	"fmt"

	"github.com/m4gshm/container/immutable"
)

func New[T any](values []T) *Slice[T] {
	elements := make([]T, len(values))
	copy(elements, values)
	return Wrap(elements)
}

func Wrap[T any](elements []T) *Slice[T] {
	return &Slice[T]{Slice: immutable.Wrap(elements), elements:  &elements}
}

type Slice[T any] struct {
	*immutable.Slice[T]
	elements   *[]T
	changeMark int32
}

var _ Vector[any] = (*Slice[any])(nil)
var _ fmt.Stringer = (*Slice[any])(nil)


func (s *Slice[T]) Add(v T) bool {
	markOnStart := s.changeMark
	*s.elements = append(*s.elements, v)
	markOnFinish := s.changeMark
	if markOnFinish != markOnStart {
		panic("concurrent Slice read and write")
	}
	return true
}


