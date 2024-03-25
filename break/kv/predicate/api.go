// Package predicate provides helpers for filtering keys or values of a map
package predicate

// Key adapts a key appliable predicate to a key\value one
func Key[V, K any](predicate func(K) (bool, error)) (out func(K, V) (bool, error)) {
	if predicate != nil {
		out = func(key K, _ V) (bool, error) { return predicate(key) }
	}
	return
}

// Value adapts a value appliable predicate to a key\value one
func Value[K, V any](predicate func(V) (bool, error)) (out func(K, V) (bool, error)) {
	if predicate != nil {
		out = func(_ K, val V) (bool, error) { return predicate(val) }
	}
	return
}
