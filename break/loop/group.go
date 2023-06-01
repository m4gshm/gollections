package loop

import (
	"github.com/m4gshm/gollections/break/kv"
)

// NewKeyValuer creates instance of the KeyValuer
func NewKeyValuer[T any, K, V any](next func() (T, bool, error), keyExtractor func(T) (K, error), valsExtractor func(T) (V, error)) KeyValuer[T, K, V] {
	return KeyValuer[T, K, V]{next: next, keyExtractor: keyExtractor, valExtractor: valsExtractor}
}

// NewMultipleKeyValuer creates instance of the MultipleKeyValuer
func NewMultipleKeyValuer[T any, K, V any](next func() (T, bool, error), keysExtractor func(T) ([]K, error), valsExtractor func(T) ([]V, error)) *MultipleKeyValuer[T, K, V] {
	return &MultipleKeyValuer[T, K, V]{next: next, keysExtractor: keysExtractor, valsExtractor: valsExtractor}
}

// KeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type KeyValuer[T, K, V any] struct {
	next         func() (T, bool, error)
	keyExtractor func(T) (K, error)
	valExtractor func(T) (V, error)
}

var _ kv.Iterator[int, string] = (*KeyValuer[any, int, string])(nil)
var _ kv.Iterator[int, string] = KeyValuer[any, int, string]{}

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv KeyValuer[T, K, V]) Track(traker func(key K, value V) error) error {
	return Track(kv.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv KeyValuer[T, K, V]) Next() (key K, value V, ok bool, err error) {
	if next := kv.next; next != nil {
		if elem, nextOk, err := next(); err != nil || !nextOk {
			return key, value, false, err
		} else if key, err = kv.keyExtractor(elem); err == nil {
			value, err = kv.valExtractor(elem)
		}
	}
	return key, value, true, err
}

// MultipleKeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type MultipleKeyValuer[T, K, V any] struct {
	next          func() (T, bool, error)
	keysExtractor func(T) ([]K, error)
	valsExtractor func(T) ([]V, error)
	keys          []K
	values        []V
	ki, vi        int
}

var _ kv.Iterator[int, string] = (*MultipleKeyValuer[any, int, string])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv *MultipleKeyValuer[T, K, V]) Track(traker func(key K, value V) error) error {
	return Track(kv.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv *MultipleKeyValuer[T, K, V]) Next() (key K, value V, ok bool, err error) {
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
				if ok {
					if kv.ki < lastKeyIndex {
						kv.ki++
					} else if kv.vi < lastValIndex {
						kv.ki = 0
						kv.vi++
					} else {
						kv.keys, kv.values = nil, nil
					}
				} else if elem, nextOk, err := next(); err != nil {
					return key, value, ok, err
				} else if nextOk {
					kv.keys, err = kv.keysExtractor(elem)
					if err == nil {
						kv.values, err = kv.valsExtractor(elem)
					}
					if err != nil {
						break
					}
					kv.ki, kv.vi = 0, 0
				} else {
					kv.keys, kv.values = nil, nil
					break
				}
			}
		}
	}
	return key, value, ok, nil
}
