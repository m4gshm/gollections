// Package group provides short aliases for functions used to group iterator elements
package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/kv/stream"
)

// Of - group.Of synonym for the iter.Group
func Of[T any, K comparable, IT c.Iterator[T]](elements IT, by func(T) K) stream.Iter[K, T, map[K][]T] {
	return iter.Group(elements, by)
}
