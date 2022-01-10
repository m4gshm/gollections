package group

import (
	"github.com/m4gshm/gollections/typ"
	"github.com/m4gshm/gollections/walk"
)

func Of[T any, K comparable, W typ.Walk[T]](elements W, by typ.Converter[T, K]) map[K][]T {
	return walk.Group(elements, by)
}
