// Package group provides short aliases for functions used to group collection elements
package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/kv/stream"
)

// Of - group.Of synonym of the c.Group
func Of[I c.Iterator[T], T any, K comparable, IT c.Iterable[I]](elements IT, by func(T) K) stream.Iter[K, T, map[K][]T] {
	return iter.Group(elements.Begin(), by)
}
