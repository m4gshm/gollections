package set

import (
	"github.com/m4gshm/container/immutable"
)

func Of[T comparable](values ...T) immutable.Set[T] {
	return immutable.NewOrderedSet(values)
}

func New[T comparable](values []T) immutable.Set[T] {
	return immutable.NewOrderedSet(values)
}