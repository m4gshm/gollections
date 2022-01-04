package set

import (
	"github.com/m4gshm/container/immutable"
	"github.com/m4gshm/container/typ"
)

func Of[T comparable](elements ...T) immutable.Set[T, typ.Iterator[T]] {
	return ToOrderedSet(elements)
}

func New[T comparable](capacity int) immutable.Set[T, typ.Iterator[T]] {
	return NewOrderedSet[T](capacity)
}