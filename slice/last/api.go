package last

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

// Of an alias of the slice.Last
func Of[T any, TS ~[]T](elements TS, by c.Predicate[T]) (T, bool) {
	return slice.Last(elements, by)
}
