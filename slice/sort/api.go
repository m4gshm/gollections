package sort

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
	"golang.org/x/exp/constraints"
)

//ByOrdered sorts elements by converting them to Ordered values and applying the operator <
func ByOrdered[T any, o constraints.Ordered, TS ~[]T](elements TS, by c.Converter[T, o]) []T {
	return slice.SortByOrdered(elements, by)
}