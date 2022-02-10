package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/mutable"
)

func Of[T any](elements ...T) *mutable.Vector[T] {
	return mutable.ToVector(elements)
}

func Empty[T any]() *mutable.Vector[T] {
	return New[T](0)
}

func New[T any](capacity int) *mutable.Vector[T] {
	return mutable.NewVector[T](capacity)
}

func Sort[t any, f constraints.Ordered](v *mutable.Vector[t], by c.Converter[t, f]) *mutable.Vector[t] {
	return v.Sort(func(e1, e2 t) bool { return by(e1) < by(e2) })
}

