package it

import (
	"github.com/m4gshm/gollections/typ"
)

func NewKeyValuer[k comparable, v any, IT typ.Iterator[v]](iter IT, keyExtractor typ.Converter[v, k]) *KeyValuer[k, v] {
	return &KeyValuer[k, v]{iter: iter, getKey: keyExtractor}
}

type KeyValuer[k comparable, v any] struct {
	iter   typ.Iterator[v]
	getKey typ.Converter[v, k]
	err    error
}

var _ typ.KVIterator[any, any] = (*KeyValuer[any, any])(nil)

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
