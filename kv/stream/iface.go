// Package stream provides a stream implementation and helper functions
package stream

import (
	"github.com/m4gshm/gollections/kv"
)

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V any, I kv.Iterator[K, V], M map[K]V | map[K][]V] interface {
	kv.Iterator[K, V]
	kv.Collection[K, V, I, M]

	HasAny(func(K, V) bool) bool
}
