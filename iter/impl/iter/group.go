package iter

import (
	"github.com/m4gshm/gollections/c"
)

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[T any, K, V any, IT any](iterator IT, next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) KeyValuer[T, K, V, IT] {
	return KeyValuer[T, K, V, IT]{iterator: iterator, next: next, key: keyExtractor, val: valExtractor}
}

// KeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type KeyValuer[T, K, V any, IT any] struct {
	iterator IT
	next     func() (T, bool)
	key      func(T) K
	val      func(T) V
}

var _ c.KVIterator[int, string] = (*KeyValuer[any, int, string, c.Iterator[any]])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv KeyValuer[T, K, V, IT]) Next() (key K, value V, ok bool) {
	if next := kv.next; next != nil {
		if elem, kvOk := next(); kvOk {
			key = kv.key(elem)
			value = kv.val(elem)
			ok = true
		}
	}
	return key, value, ok
}
