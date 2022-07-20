package ordered

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice"
)

//Sort sorts elements in place and returns them
func Sort[T constraints.Ordered, TS ~[]T](elements TS) []T {
	return slice.SortByOrdered(elements, func(o T) T { return o })
}
