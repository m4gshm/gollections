package iterator

import (
	"reflect"
)

func newMap[K comparable, V any](items map[K]V) *Map[K, V] {
	return &Map[K, V]{items: items, iter: reflect.ValueOf(items).MapRange()}
}

type Map[K comparable, V any] struct {
	items map[K]V
	iter  *reflect.MapIter
}

type KV[K comparable, V any] struct {
	key   K
	value V
}

func (k *KV[K, V]) Key() K {
	return k.key
}

func (k *KV[K, V]) Value() V {
	return k.value
}

var _ Iterator[*KV[interface{}, interface{}]] = (*Map[interface{}, interface{}])(nil)

func (s *Map[K, V]) Next() bool {
	return s.iter.Next()
}

func (s *Map[K, V]) Get() *KV[K, V] {
	key := s.iter.Key().Interface().(K)
	return &KV[K, V]{key: key, value: s.items[key]}
}
