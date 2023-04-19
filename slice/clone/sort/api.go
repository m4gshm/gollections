package sort

import (
	"sort"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"golang.org/x/exp/constraints"
)

// By makes clone of sorted elements by converting them to Ordered values and applying the operator <
func By[T any, o constraints.Ordered, TS ~[]T](elements TS, by func(T) o) TS {
	c := clone.Of(elements)
	slice.SortByOrdered(c, sort.Slice, by)
	return c
}

// ByLess makes clone and sorts elements using a function that checks if an element is smaller than the others
func ByLess[T any, TS ~[]T](elements TS, less slice.Less[T]) TS {
	c := clone.Of(elements)
	slice.Sort(c, sort.Slice, less)
	return c
}

// Of makes clone of sorted orderable elements
func Of[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return By(elements, func(o T) T { return o })
}
