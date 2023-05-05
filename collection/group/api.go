// Package group provides short aliases for functions that are used to group collection elements
package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/kv/stream"
)

// Of - group.Of synonym of the iter.Group
func Of[T any, K comparable, IT c.Iterable[T]](elements IT, by func(T) K) stream.Iter[K, T, map[K][]T] {
	return iter.Group(elements.Iter(), by)
}
