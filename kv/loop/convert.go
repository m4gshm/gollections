package loop

import (
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/loop"
)

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertIter[K, V any, K2, V2 any, C func(K, V) (K2, V2)] struct {
	next      func() (K, V, bool)
	converter C
}

var (
	_ kv.Iterator[any, any] = (*ConvertIter[any, any, any, any, func(any, any) (any, any)])(nil)
	_ kv.Iterator[any, any] = ConvertIter[any, any, any, any, func(any, any) (any, any)]{}
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i ConvertIter[K, V, K2, V2, C]) Track(traker func(key K2, value V2) error) error {
	return loop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i ConvertIter[K, V, K2, V2, C]) TrackEach(traker func(key K2, value V2)) {
	loop.TrackEach(i.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether an pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i ConvertIter[K, V, K2, V2, C]) Next() (k2 K2, v2 V2, ok bool) {
	if next := i.next; next != nil {
		if K, V, ok := next(); ok {
			k2, v2 = i.converter(K, V)
			return k2, v2, true
		}
	}
	return k2, v2, false
}

func (i ConvertIter[K, V, K2, V2, C]) Start() (ConvertIter[K, V, K2, V2, C], K2, V2, bool) {
	return startKvIt[K2, V2](i)
}
