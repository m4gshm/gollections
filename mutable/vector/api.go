package vector

import (
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
)

func Of[T any](elements ...T) mutable.Vector[T, typ.Iterator[T]] {
	return Convert(elements)
}

func New[T any](capacity int) mutable.Vector[T, typ.Iterator[T]] {
	return Create[T](capacity)
}