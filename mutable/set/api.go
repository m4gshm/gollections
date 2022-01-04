package set

import (
	"github.com/m4gshm/container/mutable"
)

func Of[T comparable](elements ...T) mutable.Set[T, mutable.Iterator[T]] {
	return ToOrderedSet(elements)
}

func New[T comparable](capacity int) mutable.Set[T, mutable.Iterator[T]] {
	return NewOrderedSet[T](capacity)
}