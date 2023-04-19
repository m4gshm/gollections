package last

import (
	"github.com/m4gshm/gollections/slice"
)

// Of an alias of the slice.Last
func Of[TS ~[]T, T any](elements TS, by func(T) bool) (T, bool) {
	return slice.Last(elements, by)
}
