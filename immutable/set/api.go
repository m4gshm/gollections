//Package set provides the unordered set container implementation
package set

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/ordered"
)

//Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) immutable.Set[T] {
	return immutable.NewSet(elements)
}

//New instantiates Set and copies elements to it.
func New[T comparable](elements []T) immutable.Set[T] {
	return immutable.NewSet(elements)
}

//Sort instantiates Set and puts sorted elements to it.
func Sort[T comparable, f constraints.Ordered](s immutable.Set[T], by c.Converter[T, f]) ordered.Set[T] {
	return s.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}
