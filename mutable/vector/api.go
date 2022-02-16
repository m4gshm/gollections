package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/mutable"
)

//Of creates the Vector with predefined elements.
func Of[T any](elements ...T) *mutable.Vector[T] {
	return mutable.ToVector(elements)
}

//Empty creates the Vector with zero capacity.
func Empty[T any]() *mutable.Vector[T] {
	return New[T](0)
}

//New creates a vector with a predefined capacity.
func New[T any](capacity int) *mutable.Vector[T] {
	return mutable.NewVector[T](capacity)
}

//Sort sorts a Vector in-place by a converter that thransforms a element to an Ordered (int, string and so on).
func Sort[T any, F constraints.Ordered](v *mutable.Vector[T], by c.Converter[T, F]) *mutable.Vector[T] {
	return v.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}
