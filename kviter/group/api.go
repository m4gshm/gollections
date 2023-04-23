package group

import "github.com/m4gshm/gollections/c"
import "github.com/m4gshm/gollections/loop/kv"

// Of collects sets of values grouped by keys obtained by passing a key/value iterator
func Of[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	return kv.Group(it.Next)
}
