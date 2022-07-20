package ordered

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice"
)

//Sort returns sorted elements
func Sort[T constraints.Ordered, TS ~[]T](elements TS) []T {
	return slice.SortByOrdered(elements, func(o T) T { return o })
}
