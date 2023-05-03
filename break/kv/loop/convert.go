package loop

import (
	"github.com/m4gshm/gollections/break/kv"
	"github.com/m4gshm/gollections/break/loop"
)

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertIter[K, V any, K2, V2 any] struct {
	next      func() (K, V, bool, error)
	converter func(K, V) (K2, V2, error)
}

var (
	_ kv.Iterator[any, any] = (*ConvertIter[any, any, any, any])(nil)
	_ kv.Iterator[any, any] = ConvertIter[any, any, any, any]{}
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i ConvertIter[K, V, K2, V2]) Track(traker func(key K2, value V2) error) error {
	return loop.Track(i.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether an pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i ConvertIter[K, V, K2, V2]) Next() (k2 K2, v2 V2, ok bool, err error) {
	if next := i.next; next != nil {
		k, v, ok, err := next()
		if err != nil || !ok {
			return k2, v2, false, err
		}
		k2, v2, err = i.converter(k, v)
		return k2, v2, err == nil, err
	}
	return k2, v2, false, nil
}
