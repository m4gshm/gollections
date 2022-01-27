//package vector provides the Vector (ordered) implementation
package vector

import (
	"github.com/m4gshm/gollections/immutable"
)

func Of[T any](elements ...T) immutable.Vector[T] {
	return Convert(elements)
}

func New[T any](elements []T) immutable.Vector[T] {
	return Convert(elements)
}
