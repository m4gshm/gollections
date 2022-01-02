package vector

import (
	"github.com/m4gshm/container/mutable"
)

func Of[T any](elements ...T) mutable.Vec[T] {
	return mutable.NewVector(elements)
}

func New[T any](elements []T) *mutable.Vector[T] {
	return mutable.NewVector(elements)
}