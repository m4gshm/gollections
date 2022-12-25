package it

import "github.com/m4gshm/gollections/c"

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[T any, K, V any, IT c.Iterator[T]](iter IT, keyExtractor c.Converter[T, K], valExtractor c.Converter[T, V]) KeyValuer[T, K, V, IT] {
	return KeyValuer[T, K, V, IT]{iter: iter, getKey: keyExtractor, getVal: valExtractor}
}

// KeyValuer is the Iterator wrapper that converts a element to a key\value pair and iterates over these pairs
type KeyValuer[T, K, V any, IT c.Iterator[T]] struct {
	iter   IT
	getKey c.Converter[T, K]
	getVal c.Converter[T, V]
}

var _ c.KVIterator[int, string] = (*KeyValuer[any, int, string, c.Iterator[any]])(nil)
var _ c.KVIterator[int, string] = KeyValuer[any, int, string, c.Iterator[any]] {}

func (s KeyValuer[T, K, V, IT]) Next() (K, V, bool) {
	elem, ok := s.iter.Next()
	if !ok {
		var k K
		var v V
		return k, v, false
	}
	k := s.getKey(elem)
	v := s.getVal(elem)
	return k, v, true
}
