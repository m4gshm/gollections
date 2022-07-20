//Package vector provides the Vector (ordered) implementation
package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
)

//Of instantiates Vector with predefined elements.
func Of[T any](elements ...T) immutable.Vector[T] {
	return immutable.NewVector(elements)
}

//New instantiates Vector and copies elements to it.
func New[T any](elements []T) immutable.Vector[T] {
	return immutable.NewVector(elements)
}

//Sort instantiates Vector and puts sorted elements to it.
func Sort[t any, f constraints.Ordered](v immutable.Vector[t], by c.Converter[t, f]) immutable.Vector[t] {
	return v.Sort(func(e1, e2 t) bool { return by(e1) < by(e2) })
}
