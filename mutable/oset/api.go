package oset

import "github.com/m4gshm/gollections/mutable/ordered"

//Of creates the Set with predefined elements.
func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.ToSet(elements)
}

//Empty creates the Set with zero capacity.
func Empty[T comparable]() *ordered.Set[T] {
	return New[T](0)
}

//New creates the Set with a predefined capacity.
func New[T comparable](capacity int) *ordered.Set[T] {
	return ordered.NewSet[T](capacity)
}
