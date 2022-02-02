package it

import "github.com/m4gshm/gollections/c"

func NewKeyValuer[k comparable, v any, IT c.Iterator[v]](iter IT, keyExtractor c.Converter[v, k]) *KeyValuer[k, v] {
	return &KeyValuer[k, v]{iter: iter, getKey: keyExtractor}
}

type KeyValuer[k comparable, v any] struct {
	iter   c.Iterator[v]
	getKey c.Converter[v, k]
	err    error
}

var _ c.KVIterator[int, any] = (*KeyValuer[int, any])(nil)

func (s *KeyValuer[k, v]) HasNext() bool {
	return s.iter.HasNext()
}

func (s *KeyValuer[k, v]) Next() (k, v) {
	val := s.iter.Next()
	return s.getKey(val), val
}
