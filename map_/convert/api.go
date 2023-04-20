package convert

// Key adapts a key converter to the key/value converter that converts only keys
func Key[V, K, KOUT any](converter func(K) KOUT) (out func(key K, val V) (KOUT, V)) {
	return func(key K, val V) (KOUT, V) { return converter(key), val }
}

// Value adapts a value converter to the key/value converter that converts only values
func Value[K, V, VOUT any](converter func(V) VOUT) (out func(key K, val V) (K, VOUT)) {
	return func(key K, val V) (K, VOUT) { return key, converter(val) }
}
