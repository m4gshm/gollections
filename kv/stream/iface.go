package stream

import (
	"github.com/m4gshm/gollections/kv"
)

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V any, M map[K]V | map[K][]V] interface {
	kv.KVIterator[K, V]
	kv.Collection[K, V, M]

	HasAny(func(K, V) bool) bool
}
