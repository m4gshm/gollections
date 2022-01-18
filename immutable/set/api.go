package set

import (
	"github.com/m4gshm/gollections/immutable"
)

func Of[T comparable](elements ...T) immutable.Set[T] {
	return Convert(elements)
}

func Empty[T comparable]() immutable.Set[T] {
	return NewImpl[T](0)
}

func New[T comparable](capacity int) immutable.Set[T] {
	return NewImpl[T](capacity)
}
