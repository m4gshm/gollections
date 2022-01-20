package set

import (
	"github.com/m4gshm/gollections/mutable"
)

func Of[T comparable](elements ...T) mutable.Set[T, mutable.Iterator[T]] {
	return Convert(elements)
}

func Empty[T comparable]() mutable.Set[T, mutable.Iterator[T]] {
	return New[T](0)
}

func New[T comparable](capacity int) mutable.Set[T, mutable.Iterator[T]] {
	return Wrap[T](make(map[T]struct{}, capacity))
}
