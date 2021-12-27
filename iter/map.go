package iter

import (
	"reflect"

	"github.com/m4gshm/container/typ"
)

func NewMap[K comparable, V any](elements map[K]V) *MapIter[K, V] {
	v := reflect.ValueOf(elements)
	return &MapIter[K, V]{elements: elements, v: v, iter: v.MapRange()}
}

type MapIter[K comparable, V any] struct {
	elements map[K]V
	iter     *reflect.MapIter
	v        reflect.Value
}

func (s *MapIter[K, V]) HasNext() bool {
	return s.iter.Next()
}

func (s *MapIter[K, V]) Get() *KV[K, V] {
	key := s.iter.Key().Interface().(K)
	return &KV[K, V]{key: key, value: s.elements[key]}
}

func (s *MapIter[K, V]) Reset() {
	s.iter.Reset(s.v)
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

func (k *KV[K, V]) Get() (K,V) {
	return k.key, k.value
}

var _ typ.Iterator[*KV[interface{}, interface{}]] = (*MapIter[interface{}, interface{}])(nil)
var _ typ.Resetable = (*MapIter[interface{}, interface{}])(nil)
