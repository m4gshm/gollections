package oset

import "github.com/m4gshm/gollections/mutable/ordered"

func Of[T comparable](elements ...T) *ordered.OrderedSet[T] {
	return ordered.Convert(elements)
}

func Empty[T comparable]() *ordered.OrderedSet[T] {
	return New[T](0)
}

func New[T comparable](capacity int) *ordered.OrderedSet[T] {
	return ordered.NewOrderedSet[T](capacity)
}
