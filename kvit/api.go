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
