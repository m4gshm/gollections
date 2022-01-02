package set

import (
	"github.com/m4gshm/container/mutable"
)

func Of[T comparable](elements ...T) mutable.Set[T] {
	return ToOrderedSet(elements)
}

func New[T comparable]() mutable.Set[T] {
	return NewOrderedSet[T]()
}