// Package filter provides helpers for filtering keys or values of a map
package filter

// Key adapts a key appliable predicate to a key\value one
func Key[V, K any](predicate func(K) bool) (out func(K, V) bool) {
	if predicate != nil {
		out = func(key K, val V) bool { return predicate(key) }
	}
	return
}

// Value adapts a value appliable predicate to a key\value one
func Value[K, V any](predicate func(V) bool) (out func(K, V) bool) {
	if predicate != nil {
		out = func(key K, val V) bool { return predicate(val) }
	}
	return
}
