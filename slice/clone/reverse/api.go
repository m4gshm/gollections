package reverse

import (
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
)

// Of clones a slice of elements and reverses the order
func Of[TS ~[]T, T any](elements TS) TS {
	return slice.Reverse(clone.Of(elements))
}
