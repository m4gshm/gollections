package clone

import "github.com/m4gshm/gollections/slice"

//Of makes new slice instance with copied elements.
func Of[T any, TS ~[]T](elements TS) []T {
	return slice.Clone(elements)
}
