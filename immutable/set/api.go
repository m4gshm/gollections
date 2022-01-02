package set

import (
	"github.com/m4gshm/container/immutable"
)

func Of[T comparable](elements ...T) immutable.Set[T] {
	return NewOrderedSet(elements)
}

func New[T comparable](elements []T) immutable.Set[T] {
	return NewOrderedSet(elements)
}