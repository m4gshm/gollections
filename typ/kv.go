package typ

func NewKV[K any, V any](key K, value V) *KV[K, V] {
	return &KV[K, V]{K: key, V: value}
}

type KV[k any, v any] struct {
	K k
	V v
}

func (k *KV[K, V]) Key() K {
	return k.K
}

func (k *KV[K, V]) Value() V {
	return k.V
}

func (k *KV[K, V]) Get() (K, V) {
	return k.K, k.V
}
