//Package oset provides the ordered set container implementation
package oset

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"

)

func Of[t comparable](elements ...t) *ordered.Set[t] {
	return ordered.NewSet(elements)
}

func New[t comparable](elements []t) *ordered.Set[t] {
	return ordered.NewSet(elements)
}

func Sort[t comparable, f constraints.Ordered](s *ordered.Set[t], by c.Converter[t, f]) *ordered.Set[t] {
	return s.Sort(func(e1, e2 t) bool { return by(e1) < by(e2) })
}