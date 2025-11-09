// Package asc provides aliases for storing slice functions.
package asc

import (
	"cmp"

	"github.com/m4gshm/gollections/slice/sort"
)

// Of is a short alias for sort.By
func Of[T any, O cmp.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return sort.By(elements, orderConverner)
}
