// Package sort provides sorting of cloned slice elements
package sort

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/sort"
)

// By makes clone of sorted elements by converting them to Ordered values and applying the operator <
func By[T any, O constraints.Ordered, TS ~[]T](elements TS, order func(T) O) TS {
	return sort.By(clone.Of(elements), order)
}

func DescBy[T any, O constraints.Ordered, TS ~[]T](elements TS, order func(T) O) TS {
	return sort.DescBy(clone.Of(elements), order)
}

// Asc sorts orderable elements ascending
func Asc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return sort.Asc(clone.Of(elements))
}

// Desc sorts orderable elements descending
func Desc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return sort.Desc(clone.Of(elements))
}
