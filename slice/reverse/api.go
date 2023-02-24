package reverse

import (
	"github.com/m4gshm/gollections/slice"
)

// Of - shortener of the slice.Reverse function
func Of[TS ~[]T, T any](elements TS) TS {
	return slice.Reverse(elements)
}
