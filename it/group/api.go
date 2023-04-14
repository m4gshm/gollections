package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it"
)

// Of - group.Of synonym for the it.Group.
func Of[T any, K comparable, IT c.Iterator[T]](elements IT, by func(T) K) c.MapPipe[K, T, map[K][]T] {
	return it.Group[T](elements, by)
}
