package filter


// Key adapts a key appliable predicate to a key\value one
func Key[K, V any](filter func(K) bool) func(K, V) bool {
	return func(key K, val V) bool { return filter(key) }
}

// Value adapts a value appliable predicate to a key\value one
func Value[K, V any](filter func(V) bool) func(K, V) bool {
	return func(key K, val V) bool { return filter(val) }
}
