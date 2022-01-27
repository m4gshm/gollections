package oset

import (
	"github.com/m4gshm/gollections/mutable"
)

func Of[T comparable](elements ...T) mutable.Set[T] {
	return ToOrderedSet(elements)
}

func Empty[T comparable]() mutable.Set[T] {
	return New[T](0)
}

func New[T comparable](capacity int) mutable.Set[T] {
	return NewOrderedSet[T](capacity)
}