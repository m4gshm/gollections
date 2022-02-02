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

func (s *KeyValuer[k, v]) Get() (k, v, error) {
	val, err := s.iter.Get()
	if err != nil {
		var key k
		var val v
		return key, val, err
	}
	return s.getKey(val), val, nil
}

func (s *KeyValuer[k, v]) Next() (k, v) {
	return NextKV[k, v](s)
}
