// Package group provides short aliases for functions used to group key/values retieved by an iterator
package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv/loop"
)

// Of collects sets of values grouped by keys obtained by passing a key/value iterator
func Of[K comparable, V any, I c.KVIterator[K, V]](elements I) map[K][]V {
	return loop.Group(elements.Next)
}
