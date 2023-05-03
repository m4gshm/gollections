// Package group provides short aliases for functions used to group collection elements
package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/walk"
)

// Of - group.Of synonym of the walk.Group.
func Of[T any, K comparable, W c.ForEachLoop[T]](elements W, by func(T) K) map[K][]T {
	return walk.Group(elements, by)
}
