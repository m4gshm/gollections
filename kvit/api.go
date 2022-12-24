package kvit

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it"
)

// Of instantiates KVIterator of predefined key\value pairs
func Of[K, V any](pairs ...c.KV[K, V]) c.KVIterator[K, V] {
	return Wrap(pairs)
}

// Wrap instantiates KVIterator using slice as the key\value pairs source
func Wrap[K, V any, P ~[]c.KV[K, V]](pairs P) c.KVIterator[K, V] {
	return it.ToKVIter[K, V](it.Wrap(pairs))
}

// Map instantiates a map with the key/values obtained by passing over a key/value iterator.
func Map[K comparable, V any](it c.KVIterator[K, V]) map[K]V {
	e := map[K]V{}
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		e[k] = v
	}
	return e
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator.
func Group[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	e := map[K][]V{}
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		e[k] = append(e[k], v)
	}
	return e
}