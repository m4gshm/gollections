package set

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/mutable/ordered"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) mutable.Set[T] {
	return mutable.ToSet(elements)
}

// Empty instantiates Set with zero capacity.
func Empty[T comparable]() mutable.Set[T] {
	return New[T](0)
}

// New instantiates Set with a predefined capacity.
func New[T comparable](capacity int) mutable.Set[T] {
	return mutable.NewSet[T](capacity)
}

// Sort sorts a Set in-place by a converter that thransforms a element to an Ordered (int, string and so on).
func Sort[T comparable, F constraints.Ordered](s mutable.Set[T], by c.Converter[T, F]) *ordered.Set[T] {
	return s.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}
