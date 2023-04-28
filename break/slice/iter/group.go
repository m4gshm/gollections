package iter

import (
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/iter"
)

// KeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type KeyValuer[T, K, V any] struct {
	iter         iter.SliceIter[T]
	keyExtractor func(T) (K, error)
	valExtractor func(T) (V, error)
}

var _ c.KVIteratorBreakable[int, string] = (*KeyValuer[any, int, string])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv KeyValuer[T, K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(kv.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv *KeyValuer[T, K, V]) Next() (key K, value V, ok bool, err error) {
	next := kv.iter.Next
	if elem, nextOk := next(); nextOk {
		key, err = kv.keyExtractor(elem)
		if err == nil {
			value, err = kv.valExtractor(elem)
		}
		ok = err == nil
	}
	return key, value, ok, err
}

// MultipleKeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type MultipleKeyValuer[T, K, V any] struct {
	iter          iter.SliceIter[T]
	keysExtractor func(T) ([]K, error)
	valsExtractor func(T) ([]V, error)
	keys          []K
	values        []V
	ki, vi        int
}

var _ c.KVIteratorBreakable[int, string] = (*MultipleKeyValuer[any, int, string])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv *MultipleKeyValuer[T, K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(kv.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv *MultipleKeyValuer[T, K, V]) Next() (key K, value V, ok bool, err error) {
	if kv != nil {
		next := kv.iter.Next
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
			} else if elem, nextOk := next(); nextOk {
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
	return key, value, ok, err
}
