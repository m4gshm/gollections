// Package stablesort provides stable sorting in place slice elements
package stablesort

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/convert/as"
	"github.com/m4gshm/gollections/slice"
)

// By sorts elements in ascending order, using the orderConverner function to retrieve a value of type Ordered.
func By[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverter func(T) O) TS {
	return slice.StableSortAsc(elements, orderConverter)
}

// Asc sorts orderable elements ascending
func Asc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return slice.StableSortAsc(elements, as.Is[T])
}

// Desc sorts orderable elements descending
func Desc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return slice.StableSortDesc(elements, as.Is[T])
}
