// Package stream provides a stream implementation and helper functions
package stream

import (
	"github.com/m4gshm/gollections/kv/collection"
)

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V any, M map[K]V | map[K][]V] interface {
	collection.Iterator[K, V]
	collection.Collection[K, V, M]

	HasAny(func(K, V) bool) bool
}
