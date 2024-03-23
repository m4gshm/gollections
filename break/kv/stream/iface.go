package stream

import "github.com/m4gshm/gollections/break/kv/loop"

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V any, Map map[K]V | map[K][]V] interface {
	Loop() loop.Loop[K, V]

	Map() (Map, error)

	Reduce(merger func(K, K, V, V) (K, V, error)) (K, V, error)
	HasAny(predicate func(K, V) (bool, error)) (bool, error)
}
