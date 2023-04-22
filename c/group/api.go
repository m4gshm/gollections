package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
)

// Of - group.Of synonym of the c.Group
func Of[T any, K comparable, IT c.Iterable[T]](elements IT, by func(T) K) c.MapTransform[K, T, map[K][]T] {
	return iter.Group(elements.Begin(), by)
}
