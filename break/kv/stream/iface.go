package stream

import "github.com/m4gshm/gollections/break/kv"

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V any, Map map[K]V | map[K][]V] interface {
	kv.Iterator[K, V]
	Loop() func() (K, V, bool, error)

	// Filter(predicate func(K, V) bool) Stream[K, V, I, Map]
	// FilterKey(predicate func(K) bool) Stream[K, V, I, Map]
	// FilterValue(predicate func(V) bool) Stream[K, V, I, Map]

	// Filt(predicate func(K, V) (bool, error)) Stream[K, V, I, Map]
	// FiltKey(predicate func(K) (bool, error)) Stream[K, V, I, Map]
	// FiltValue(predicate func(V) (bool, error)) Stream[K, V, I, Map]

	Map() (Map, error)

	Reduce(merger func(K, K, V, V) (K, V, error)) (K, V, error)
	HasAny(predicate func(K, V) (bool, error)) (bool, error)
}
