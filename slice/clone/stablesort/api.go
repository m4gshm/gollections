// Package stablesort provides stable sorting of cloned slice elements
package stablesort

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
)

// By makes clone of stable sorted elements by converting them to Ordered values and applying the operator <
func By[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverter func(T) O) TS {
	return slice.StableSortAsc(clone.Of(elements), orderConverter)
}

// ByLess makes clone and atanle sorts elements using a function that checks if an element is smaller than the others
func ByLess[T any, TS ~[]T](elements TS, comparer func(T, T) int) TS {
	return slice.StableSort(clone.Of(elements), comparer)
}

// Asc makes clone of stable sorted orderable elements
func Asc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return By(elements, func(o T) T { return o })
}
