package it

import "github.com/m4gshm/gollections/c"

type ConvertFit[From, To any] struct {
	Iter    c.Iterator[From]
	By      c.Converter[From, To]
	Fit     c.Predicate[From]
	current To
}

var _ c.Iterator[any] = (*ConvertFit[any, any])(nil)

func (s *ConvertFit[From, To]) HasNext() bool {
	if v, ok := nextFiltered(s.Iter, s.Fit); ok {
		s.current = s.By(v)
		return true
	}
	return false
}

func (s *ConvertFit[From, To]) Next() To {
	return s.current
}

type Convert[From, To any, IT c.Iterator[From], C c.Converter[From, To]] struct {
	Iter IT
	By   C
}

var _ c.Iterator[any] = (*Convert[any, any, c.Iterator[any], c.Converter[any, any]])(nil)

func (s *Convert[From, To, IT, C]) HasNext() bool {
	return s.Iter.HasNext()
}

func (s *Convert[From, To, IT, C]) Next() To {
	return s.By(s.Iter.Next())
}

type ConvertKV[k, v any, IT c.KVIterator[k, v], k2, v2 any, C c.BiConverter[k, v, k2, v2]] struct {
	Iter IT
	By   C
}

var _ c.KVIterator[any, any] = (*ConvertKV[any, any, c.KVIterator[any, any], any, any, c.BiConverter[any, any, any, any]])(nil)

func (s *ConvertKV[k, v, IT, k1, v2, C]) HasNext() bool {
	return s.Iter.HasNext()
}

func (s *ConvertKV[k, v, IT, k2, v2, C]) Next() (k2, v2) {
	return s.By(s.Iter.Next())
}
