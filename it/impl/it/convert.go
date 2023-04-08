package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/predicate"
)

// ConvertFitIter is the Converter with elements filtering.
type ConvertFitIter[From, To any, IT c.Iterator[From]] struct {
	iter IT
	by   c.Converter[From, To]
	fit  predicate.Predicate[From]
}

var (
	_ c.Iterator[any] = ConvertFitIter[any, any, c.Iterator[any]]{}
	_ c.Iterator[any] = (*ConvertFitIter[any, any, c.Iterator[any]])(nil)
)

func (s ConvertFitIter[From, To, IT]) Next() (To, bool) {
	if V, ok := nextFiltered(s.iter, s.fit); ok {
		return s.by(V), true
	}
	var no To
	return no, false
}

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable elements.
type ConvertIter[From, To any, IT c.Iterator[From], C c.Converter[From, To]] struct {
	iter IT
	by   C
}

var (
	_ c.Iterator[any] = ConvertIter[any, any, c.Iterator[any], c.Converter[any, any]]{}
	_ c.Iterator[any] = (*ConvertIter[any, any, c.Iterator[any], c.Converter[any, any]])(nil)
)

func (s ConvertIter[From, To, IT, C]) Next() (To, bool) {
	if v, ok := s.iter.Next(); ok {
		return s.by(v), true
	}
	var no To
	return no, false
}

// ConvertKVIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertKVIter[K, V any, IT c.KVIterator[K, V], K2, V2 any, C c.BiConverter[K, V, K2, V2]] struct {
	iter IT
	by   C
}

var (
	_ c.KVIterator[any, any] = ConvertKVIter[any, any, c.KVIterator[any, any], any, any, c.BiConverter[any, any, any, any]]{}
	_ c.KVIterator[any, any] = (*ConvertKVIter[any, any, c.KVIterator[any, any], any, any, c.BiConverter[any, any, any, any]])(nil)
)

func (s ConvertKVIter[K, V, IT, K2, V2, C]) Next() (K2, V2, bool) {
	if K, V, ok := s.iter.Next(); ok {
		k2, v2 := s.by(K, V)
		return k2, v2, true
	}
	var k2 K2
	var v2 V2
	return k2, v2, false
}
