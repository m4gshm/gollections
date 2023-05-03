// Package first provides short aliases for slice functions for retrieving a first element
package first

import (
	"github.com/m4gshm/gollections/slice"
)

// Of an alias of the slice.First
func Of[TS ~[]T, T any](elements TS, filter func(T) bool) (T, bool) {
	return slice.First(elements, filter)
}
