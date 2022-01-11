package set

import (
	"github.com/m4gshm/gollections/immutable"
)

func Of[T comparable](elements ...T) immutable.Set[T] {
	return ToOrderedSet(elements)
}

func New[T comparable](capacity int) immutable.Set[T] {
	return NewOrderedSet[T](capacity)
}
