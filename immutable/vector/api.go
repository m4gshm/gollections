package vector

import (
	"github.com/m4gshm/container/immutable"
)

func Of[T any](elements ...T) immutable.Vec[T] {
	return immutable.NewVector(elements)
}

func New[T any](elements []T) *immutable.Vector[T] {
	return immutable.NewVector(elements)
}