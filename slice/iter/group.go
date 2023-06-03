package iter

import (
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// KeyValuerIter is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type KeyValuerIter[T, K, V any] struct {
	iter         slice.Iter[T]
	keyExtractor func(T) K
	valExtractor func(T) V
}

var _ kv.Iterator[int, string] = (*KeyValuerIter[any, int, string])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv KeyValuerIter[T, K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(kv.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (kv *KeyValuerIter[T, K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(kv.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv *KeyValuerIter[T, K, V]) Next() (key K, value V, ok bool) {
	next := kv.iter.Next
	if elem, nextOk := next(); nextOk {
		key = kv.keyExtractor(elem)
		value = kv.valExtractor(elem)
		ok = true
	}
	return key, value, ok
}

func (kv *KeyValuerIter[T, K, V]) Start() (*KeyValuerIter[T, K, V], K, V, bool) {
	return startKvIt[K, V](kv)
}

// MultipleKeyValuer is the Iterator wrapper that converts an element to a key\value pair and iterates over these pairs
type MultipleKeyValuer[T, K, V any] struct {
	iter          slice.Iter[T]
	keysExtractor func(T) []K
	valsExtractor func(T) []V
	keys          []K
	values        []V
	ki, vi        int
}

var _ kv.Iterator[int, string] = (*MultipleKeyValuer[any, int, string])(nil)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv *MultipleKeyValuer[T, K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(kv.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (kv *MultipleKeyValuer[T, K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(kv.Next, traker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (kv *MultipleKeyValuer[T, K, V]) Next() (key K, value V, ok bool) {
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

			if kv.ki < lastKeyIndex {
				kv.ki++
			} else if kv.vi < lastValIndex {
				kv.ki = 0
				kv.vi++
			} else if elem, nextOk := next(); nextOk {
				kv.keys = kv.keysExtractor(elem)
				kv.values = kv.valsExtractor(elem)
				kv.ki, kv.vi = 0, 0
			} else {
				kv.keys, kv.values = nil, nil
				break
			}
		}
	}
	return key, value, ok
}

func (kv *MultipleKeyValuer[T, K, V]) Start() (*MultipleKeyValuer[T, K, V], K, V, bool) {
	return startKvIt[K, V](kv)
}
