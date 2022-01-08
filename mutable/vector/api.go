package vector

import (
	"github.com/m4gshm/container/mutable"
)

func Of[T any](elements ...T) mutable.Vector[T, mutable.Iterator[T]] {
	return Convert(elements)
}

func New[T any](capacity int) mutable.Vector[T, mutable.Iterator[T]] {
	return Create[T](capacity)
}
