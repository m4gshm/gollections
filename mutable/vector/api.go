package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it"
	"github.com/m4gshm/gollections/mutable"
)

// Of instantiates Vector with predefined elements.
func Of[T any](elements ...T) *mutable.Vector[T] {
	return mutable.NewVector(elements)
}

// Empty instantiates Vector with zero capacity.
func Empty[T any]() *mutable.Vector[T] {
	return NewCap[T](0)
}

// NewCap creates a vector with a predefined capacity.
func NewCap[T any](capacity int) *mutable.Vector[T] {
	return mutable.NewVectorCap[T](capacity)
}

// From creates a Vector instance with elements obtained by passing an iterator.
func From[T any](elements c.Iterator[T]) *mutable.Vector[T] {
	return mutable.WrapVector(it.ToSlice(elements))
}

// Sort sorts a Vector in-place by a converter that thransforms a element to an Ordered (int, string and so on).
func Sort[T any, F constraints.Ordered](v *mutable.Vector[T], by func(T) F) *mutable.Vector[T] {
	return v.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}
