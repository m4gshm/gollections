package group

import (
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/typ"
)

//Of - group.Of synonym of the slice.Group.
func Of[T any, K comparable](elements []T, by typ.Converter[T, K]) map[K][]T {
	return slice.Group(elements, by)
}
