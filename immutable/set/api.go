//Package set provides the unordered set container implementation
package set

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/ordered"
)

func Of[t comparable](elements ...t) *immutable.Set[t] {
	return immutable.NewSet(elements)
}

func New[t comparable](elements []t) *immutable.Set[t] {
	return immutable.NewSet(elements)
}

func Sort[t comparable, f constraints.Ordered](s *immutable.Set[t], by c.Converter[t, f]) *ordered.Set[t] {
	return s.Sort(func(e1, e2 t) bool { return by(e1) < by(e2) })
}
