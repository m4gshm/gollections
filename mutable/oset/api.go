package oset

import "github.com/m4gshm/gollections/mutable/ordered"

//Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.ToSet(elements)
}

//Empty instantiates Set with zero capacity.
func Empty[T comparable]() *ordered.Set[T] {
	return New[T](0)
}

//New instantiates Set with a predefined capacity.
func New[T comparable](capacity int) *ordered.Set[T] {
	return ordered.NewSet[T](capacity)
}
