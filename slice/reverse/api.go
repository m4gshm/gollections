package reverse

import (
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
)

// Of makes clone of elements in reverse order
func Of[T any, TS ~[]T](elements TS) []T {
	return slice.Reverse(clone.Of(elements))
}
