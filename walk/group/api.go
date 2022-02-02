package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/walk"
)

//Of - group.Of synonym of the walk.Group.
func Of[T any, K comparable, W c.WalkEach[T]](elements W, by c.Converter[T, K]) map[K][]T {
	return walk.Group(elements, by)
}
