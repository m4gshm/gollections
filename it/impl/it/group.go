package it

import "github.com/m4gshm/gollections/c"

func NewKeyValuer[K comparable, V any, IT c.Iterator[V]](iter IT, keyExtractor c.Converter[V, K]) *KeyValuer[K, V, IT] {
	return &KeyValuer[K, V, IT]{iter: iter, getKey: keyExtractor}
}

type KeyValuer[K comparable, V any, IT c.Iterator[V]] struct {
	iter   IT
	getKey c.Converter[V, K]
	err    error
}

var _ c.KVIterator[int, any] = (*KeyValuer[int, any, c.Iterator[any]])(nil)

func (s *KeyValuer[K, V, IT]) Next() (K, V, bool) {
	v, ok := s.iter.Next()
	if !ok {
		var k K
		var v V
		return k, v, false
	}
	k := s.getKey(v)
	return k, v, true
}

func (s *KeyValuer[K, V, IT]) Cap() int {
	return s.iter.Cap()
}
