package stablesort

import (
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	"golang.org/x/exp/constraints"
)

// By makes clone of stable sorted elements by converting them to Ordered values and applying the operator <
func By[T any, o constraints.Ordered, TS ~[]T](elements TS, by c.Converter[T, o]) TS {
	c := clone.Of(elements)
	slice.SortByOrdered(c, sort.SliceStable, by)
	return c
}

// Of makes clone of stable sorted orderable elements
func Of[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return By(elements, func(o T) T { return o })
}
