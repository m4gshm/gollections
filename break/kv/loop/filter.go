package loop

import (
	"github.com/m4gshm/gollections/break/kv"
	"github.com/m4gshm/gollections/break/loop"
)

// FiltKVIter is the KVIterator wrapper that provides filtering of key/value elements by a Predicate.
type FiltKVIter[K, V any] struct {
	next   func() (K, V, bool, error)
	filter func(K, V) (bool, error)
}

var (
	_ kv.Iterator[any, any] = (*FiltKVIter[any, any])(nil)
	_ kv.Iterator[any, any] = FiltKVIter[any, any]{}
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f FiltKVIter[K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(f.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f FiltKVIter[K, V]) Next() (key K, value V, ok bool, err error) {
	if !(f.next == nil || f.filter == nil) {
		key, value, ok, err = nextFiltered(f.next, f.filter)
	}
	return key, value, ok, err
}

func nextFiltered[K any, V any](next func() (K, V, bool, error), filter func(K, V) (bool, error)) (key K, val V, filtered bool, err error) {
	return Firstt(next, filter)
}
