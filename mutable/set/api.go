package set

import (
	"github.com/m4gshm/container/mutable"
)

func Of[T comparable](values ...T) mutable.Set[T] {
	return mutable.ToOrderedSet(values)
}

func New[T comparable]() mutable.Set[T] {
	return mutable.NewOrderedSet[T]()
}