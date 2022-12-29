package first

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

// Of an alias of the slice.First
func Of[TS ~[]T, T any](elements TS, by c.Predicate[T]) (T, bool) {
	return slice.First(elements, by)
}
