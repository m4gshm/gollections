package stream

import "github.com/m4gshm/gollections/c"

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V any, M map[K]V | map[K][]V] interface {
	c.KVIterator[K, V]
	c.KVCollection[K, V, M]
}
