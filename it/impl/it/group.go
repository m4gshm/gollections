package it

import "github.com/m4gshm/gollections/c"

// NewKeyValuer creates instance of the KeyValuer.
func NewKeyValuer[K comparable, V any, IT c.Iterator[V]](iter IT, keyExtractor c.Converter[V, K]) *KeyValuer[K, V, IT] {
	return &KeyValuer[K, V, IT]{iter: iter, getKey: keyExtractor}
}

// KeyValuer is the Iterator wrapper that converts a element to a key and iterates over the key/element pairs.
type KeyValuer[K comparable, V any, IT c.Iterator[V]] struct {
	iter   IT
	getKey c.Converter[V, K]
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
