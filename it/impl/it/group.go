package it

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/typ"
)

func NewGrouper[k comparable, v any, IT typ.Iterator[v]](iter IT, keyExtractor typ.Converter[v, k]) *Grouper[k, v] {
	return &Grouper[k, v]{iter: iter, getKey: keyExtractor}
}

type Grouper[k comparable, v any] struct {
	iter   typ.Iterator[v]
	getKey typ.Converter[v, k]
	err    error
}

var _ typ.Iterator[*typ.KV[any, any]] = (*Grouper[any, any])(nil)

func (s *Grouper[k, v]) HasNext() bool {
	return s.iter.HasNext()
}

func (s *Grouper[k, v]) Next() (*typ.KV[k, v], error) {
	val, err := s.iter.Next()
	if err != nil {
		return nil, err
	}
	return K.V(s.getKey(val), val), nil
}
