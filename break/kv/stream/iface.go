package stream

import "github.com/m4gshm/gollections/c"

// Stream is map or key/value stream of elements in transformation state.
type Stream[K comparable, V, I any, Map map[K]V | map[K][]V] interface {
	c.KVIteratorBreakable[K, V]
	Iter() I

	// Filter(predicate func(K, V) bool) KVStreamBreakable[K, V, Map]
	// FilterKey(predicate func(K) bool) KVStreamBreakable[K, V, Map]
	// FilterValue(predicate func(V) bool) KVStreamBreakable[K, V, Map]

	// Filt(predicate func(K, V) (bool, error)) KVStreamBreakable[K, V, Map]
	// FiltKey(predicate func(K) (bool, error)) KVStreamBreakable[K, V, Map]
	// FiltValue(predicate func(V) (bool, error)) KVStreamBreakable[K, V, Map]

	Map() (Map, error)

	Reduce(merger func(K, V, K, V) (K, V, error)) (K, V, error)
	HasAny(predicate func(K, V) (bool, error)) (bool, error)
}
