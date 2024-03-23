// Package group provides short aliases for functions that are used to group collection elements
package group

import (
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/convert/as"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
)

// Group groups elements by keys into a map
func Of[T any, K comparable, IT collection.Iterable[T]](elements IT, by func(T) K) stream.Iter[K, T, map[K][]T] {
	return stream.New( loop.KeyValue[T, K, T](elements.Loop(), by, as.Is), kvloop.Group)
}
