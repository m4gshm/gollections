package sort

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/ordered"
	"golang.org/x/exp/constraints"
)

//ByOrdered - synony of the slice.SortByOrdered
func ByOrdered[T any, o constraints.Ordered, TS ~[]T](elements TS, by c.Converter[T, o]) []T {
	return slice.SortByOrdered(elements, by)
}


//Of makes clone of sorted elements by the ordered.Sort
func Of[T constraints.Ordered, TS ~[]T](elements TS) []T {
	return ordered.Sort(clone.Of(elements))
}