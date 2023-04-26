package loop

import (
	"github.com/m4gshm/gollections/c"
)

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[T any, K, V any](next func() (T, bool), keyProducer func(T) K, valsProducer func(T) V) KeyValuer[T, K, V] {
	return KeyValuer[T, K, V]{next: next, keyProducer: keyProducer, valProducer: valsProducer}
}

// NewMultipleKeyValuer creates instance of the MultipleKeyValuer
func NewMultipleKeyValuer[T any, K, V any](next func() (T, bool), keysProducer func(T) []K, valsProducer func(T) []V) *MultipleKeyValuer[T, K, V] {
	return &MultipleKeyValuer[T, K, V]{next: next, keysProducer: keysProducer, valsProducer: valsProducer}
}

// KeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type KeyValuer[T, K, V any] struct {
	next        func() (T, bool)
	keyProducer func(T) K
	valProducer func(T) V
}

var _ c.KVIterator[int, string] = (*KeyValuer[any, int, string])(nil)
var _ c.KVIterator[int, string] = KeyValuer[any, int, string]{}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv KeyValuer[T, K, V]) Next() (key K, value V, ok bool) {
	if next := kv.next; next != nil {
		if elem, nextOk := next(); nextOk {
			key = kv.keyProducer(elem)
			value = kv.valProducer(elem)
			ok = true
		}
	}
	return key, value, ok
}

// MultipleKeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type MultipleKeyValuer[T, K, V any] struct {
	next         func() (T, bool)
	keysProducer func(T) []K
	valsProducer func(T) []V
	keys         []K
	values       []V
	ki, vi       int
}

var _ c.KVIterator[int, string] = (*MultipleKeyValuer[any, int, string])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv *MultipleKeyValuer[T, K, V]) Next() (key K, value V, ok bool) {
	if kv != nil {
		if next := kv.next; next != nil {
			for !ok {
				var (
					keys, values               = kv.keys, kv.values
					keysLen, valuesLen         = len(keys), len(values)
					lastKeyIndex, lastValIndex = keysLen - 1, valuesLen - 1
				)

				if keysLen > 0 && kv.ki >= 0 && kv.ki <= lastKeyIndex {
					key = keys[kv.ki]
					ok = true
				}
				if valuesLen > 0 && kv.vi >= 0 && kv.vi <= lastValIndex {
					value = values[kv.vi]
					ok = true
				}

				if kv.ki < lastKeyIndex {
					kv.ki++
				} else if kv.vi < lastValIndex {
					kv.ki = 0
					kv.vi++
				} else if elem, nextOk := next(); nextOk {
					kv.keys = kv.keysProducer(elem)
					kv.values = kv.valsProducer(elem)
					kv.ki, kv.vi = 0, 0
				} else {
					kv.keys, kv.values = nil, nil
					break
				}
			}
		}
	}
	return key, value, ok
}
