package vector

import (
	"github.com/m4gshm/container/immutable"
)

func Of[T any](elements ...T) immutable.Vector[T] {
	return Convert(elements)
}

func New[T any](elements []T) immutable.Vector[T] {
	return Convert(elements)
}