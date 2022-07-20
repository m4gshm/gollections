package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

//Of - group.Of synonym of the slice.Group
func Of[T any, K comparable, TS ~[]T](elements TS, by c.Converter[T, K]) map[K][]T {
	return slice.Group(elements, by)
}
