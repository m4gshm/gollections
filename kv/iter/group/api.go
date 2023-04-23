package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv/loop"
)

// Of collects sets of values grouped by keys obtained by passing a key/value iterator
func Of[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	return loop.Group(it.Next)
}
