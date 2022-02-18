package it

import "github.com/m4gshm/gollections/c"

type ConvertFit[From, To any, IT c.Iterator[From]] struct {
	Iter    IT
	By      c.Converter[From, To]
	Fit     c.Predicate[From]
	current To
}

var _ c.Iterator[any] = (*ConvertFit[any, any, c.Iterator[any]])(nil)

func (s *ConvertFit[From, To, IT]) GetNext() (To, bool) {
	if v, ok := nextFiltered(s.Iter, s.Fit); ok {
		return s.By(v), true
	}
	var no To
	return no, false
}

func (s *ConvertFit[From, To, IT]) Next() To {
	return s.current
}

type Convert[From, To any, IT c.Iterator[From], C c.Converter[From, To]] struct {
	Iter IT
	By   C
}

var _ c.Iterator[any] = (*Convert[any, any, c.Iterator[any], c.Converter[any, any]])(nil)

func (s *Convert[From, To, IT, C]) GetNext() (To, bool) {
	if v, ok := s.Iter.GetNext(); ok {
		return s.By(v), true
	}
	var no To
	return no, false
}

type ConvertKV[k, v any, IT c.KVIterator[k, v], k2, v2 any, C c.BiConverter[k, v, k2, v2]] struct {
	Iter IT
	By   C
}

var _ c.KVIterator[any, any] = (*ConvertKV[any, any, c.KVIterator[any, any], any, any, c.BiConverter[any, any, any, any]])(nil)

func (s *ConvertKV[K, V, IT, K2, V2, C]) GetNext() (K2, V2, bool) {
	k, v, ok := s.Iter.GetNext()
	if ok {
		k2, v2 := s.By(k, v)
		return k2, v2, true
	}
	var k2 K2
	var v2 V2
	return k2, v2, false
}
