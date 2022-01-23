package vector

import (
	"github.com/m4gshm/gollections/mutable"
)

func Of[T any](elements ...T) mutable.Vector[T] {
	return Convert(elements)
}

func Empty[T any]() mutable.Vector[T] {
	return New[T](0)
}

func New[T any](capacity int) mutable.Vector[T] {
	return Create[T](capacity)
}
