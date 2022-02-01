package oset

import "github.com/m4gshm/gollections/mutable/ordered"

func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.ToSet(elements)
}

func Empty[T comparable]() *ordered.Set[T] {
	return New[T](0)
}

func New[T comparable](capacity int) *ordered.Set[T] {
	return ordered.NewSet[T](capacity)
}

func Convert[T comparable](elements []T) *ordered.Set[T] {
	return ordered.ToSet(elements)
}
