package group

import "github.com/m4gshm/gollections/c"

// Of collects sets of values grouped by keys obtained by passing a key/value iterator
func Of[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	e := map[K][]V{}
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		e[k] = append(e[k], v)
	}
	return e
}
