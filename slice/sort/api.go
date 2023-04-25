// Package sort provides sorting in place slice elements
package sort

import (
	"sort"

	"github.com/m4gshm/gollections/slice"
	"golang.org/x/exp/constraints"
)

// By sorts elements in place by converting them to Ordered values and applying the operator <
func By[T any, o constraints.Ordered, TS ~[]T](elements TS, by func(T) o) TS {
	return slice.SortByOrdered(elements, sort.Slice, by)

}

// ByLess sorts elements in place using a function that checks if an element is smaller than the others
func ByLess[T any, TS ~[]T](elements TS, less slice.Less[T]) TS {
	return slice.Sort(elements, sort.Slice, less)
}

// Of sorts orderable elements
func Of[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return By(elements, func(o T) T { return o })
}
