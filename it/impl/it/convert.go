package it

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertFitIter is the Converter with elements filtering.
type ConvertFitIter[From, To, IT any] struct {
	iter   IT
	next   func() (From, bool)
	by     func(From) To
	filter func(From) bool
}

var (
	_ c.Iterator[any] = ConvertFitIter[any, any, any]{}
	_ c.Iterator[any] = (*ConvertFitIter[any, any, any])(nil)
)

func (s ConvertFitIter[From, To, IT]) Next() (To, bool) {
	if V, ok := nextFiltered(s.next, s.filter); ok {
		return s.by(V), true
	}
	var no To
	return no, false
}

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable elements.
type ConvertIter[From, To any, IT any] struct {
	iter IT
	next func() (From, bool)
	by   func(From) To
}

var (
	_ c.Iterator[any] = ConvertIter[any, any, any]{}
	_ c.Iterator[any] = (*ConvertIter[any, any, any])(nil)
)

func (s ConvertIter[From, To, IT]) Next() (To, bool) {
	if v, ok := s.next(); ok {
		return s.by(v), true
	}
	var no To
	return no, false
}

// ConvertKVIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertKVIter[K, V any, IT c.KVIterator[K, V], K2, V2 any, C func(K, V) (K2, V2)] struct {
	iter IT
	by   C
}

var (
	_ c.KVIterator[any, any] = ConvertKVIter[any, any, c.KVIterator[any, any], any, any, func(any, any) (any, any)]{}
	_ c.KVIterator[any, any] = (*ConvertKVIter[any, any, c.KVIterator[any, any], any, any, func(any, any) (any, any)])(nil)
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
