package it

import (
	"github.com/m4gshm/gollections/c"
)

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[T any, K, V any, IT any](iter IT, next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) KeyValuer[T, K, V, IT] {
	return KeyValuer[T, K, V, IT]{iter: iter, next: next, key: keyExtractor, val: valExtractor}
}

// KeyValuer is the Iterator wrapper that converts a element to a key\value pair and iterates over these pairs
type KeyValuer[T, K, V any, IT any] struct {
	iter IT
	next func() (T, bool)
	key  func(T) K
	val  func(T) V
}

var _ c.KVIterator[int, string] = (*KeyValuer[any, int, string, c.Iterator[any]])(nil)
var _ c.KVIterator[int, string] = KeyValuer[any, int, string, c.Iterator[any]]{}

func (s KeyValuer[T, K, V, IT]) Next() (K, V, bool) {
	elem, ok := s.next()
	if !ok {
		var k K
		var v V
		return k, v, false
	}
	k := s.key(elem)
	v := s.val(elem)
	return k, v, true
}
