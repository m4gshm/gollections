package c

// NewKV creates a key/value pair holder.
func NewKV[K any, V any](key K, value V) KV[K, V] {
	return KV[K, V]{K: key, V: value}
}

// KV is the simplest implementation of a key/value pair.
type KV[k any, v any] struct {
	K k
	V v
}

// Key returns the key
func (k KV[K, V]) Key() K {
	return k.K
}

// Value returns the value
func (k KV[K, V]) Value() V {
	return k.V
}

// Get returns the key/value pair
func (k KV[K, V]) Get() (K, V) {
	return k.K, k.V
}
