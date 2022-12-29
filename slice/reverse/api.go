package reverse

import (
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
)

// Of makes clone of elements in reverse order
func Of[TS ~[]T, T any](elements TS) TS {
	return slice.Reverse(clone.Of(elements))
}
