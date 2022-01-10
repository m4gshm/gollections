package vector

import (
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/typ"
)

func Of[T any](elements ...T) immutable.Vector[T, typ.Iterator[T]] {
	return Convert(elements)
}

func New[T any](elements []T) immutable.Vector[T, typ.Iterator[T]] {
	return Convert(elements)
}
