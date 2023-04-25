// Package first provides short aliases for loop functions for retrieving a first element
package first

import (
	"github.com/m4gshm/gollections/loop"
)

// Of an alias of the loop.First
func Of[T any](next func() (T, bool), predicate func(T) bool) (T, bool) {
	return loop.First(next, predicate)
}
