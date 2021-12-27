package typ


func NewKV[K any, V any](key K, value V) *KV[K, V] {
	return &KV[K, V]{key: key, value: value}
}

type KV[K any, V any] struct {
	key   K
	value V
}

func (k *KV[K, V]) Key() K {
	return k.key
}

func (k *KV[K, V]) Value() V {
	return k.value
}
