// Package convert provides key, value convert adapters
package convert

// Key adapts a key converter to the key/value converter that converts only keys
func Key[V, K, KOUT any](converter func(K) (KOUT, error)) (out func(key K, val V) (KOUT, V, error)) {
	return func(key K, val V) (KOUT, V, error) {
		k, err := converter(key)
		return k, val, err
	}
}

// Value adapts a value converter to the key/value converter that converts only values
func Value[K, V, VOUT any](converter func(V) (VOUT, error)) (out func(key K, val V) (K, VOUT, error)) {
	return func(key K, val V) (K, VOUT, error) {
		v, err := converter(val)
		return key, v, err
	}
}
