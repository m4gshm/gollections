package iterator

import (
	"reflect"
)

func wrapMap[K comparable, V any](items map[K]V) *MapIter[K, V] {
	return &MapIter[K, V]{items: items, iter: reflect.ValueOf(items).MapRange()}
}

type MapIter[K comparable, V any] struct {
	items map[K]V
	iter  *reflect.MapIter
}

func (s *MapIter[K, V]) Next() bool {
	return s.iter.Next()
}

func (s *MapIter[K, V]) Get() *KV[K, V] {
	key := s.iter.Key().Interface().(K)
	return &KV[K, V]{key: key, value: s.items[key]}
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

var _ Iterator[*KV[interface{}, interface{}]] = (*MapIter[interface{}, interface{}])(nil)
