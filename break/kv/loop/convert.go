package loop

import (
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
)

// ConvertKVIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertKVIter[K, V any, K2, V2 any] struct {
	next      func() (K, V, bool, error)
	converter func(K, V) (K2, V2, error)
}

var (
	_ c.KVIteratorBreakable[any, any] = (*ConvertKVIter[any, any, any, any])(nil)
	_ c.KVIteratorBreakable[any, any] = ConvertKVIter[any, any, any, any]{}
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f ConvertKVIter[K, V, K2, V2]) Track(traker func(key K2, value V2) error) error {
	return loop.Track(f.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether an pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertKVIter[K, V, K2, V2]) Next() (k2 K2, v2 V2, ok bool, err error) {
	if next := c.next; next != nil {
		if k, v, ok, err := next(); err != nil || !ok {
			return k2, v2, false, err
		} else {
			k2, v2, err = c.converter(k, v)
			return k2, v2, err == nil, err
		}
	}
	return k2, v2, false, nil
}
