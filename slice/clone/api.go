package clone

import "github.com/m4gshm/gollections/slice"

// Of - synonym of the slice.Clone
func Of[T any, TS ~[]T](elements TS) []T {
	return slice.Clone(elements)
}
