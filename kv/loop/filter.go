package loop

import (
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/loop"
)

// FilterIter is the KVIterator wrapper that provides filtering of key/value elements by a Predicate.
type FilterIter[K, V any] struct {
	next   func() (K, V, bool)
	filter func(K, V) bool
}

var (
	_ kv.Iterator[any, any] = (*FilterIter[any, any])(nil)
	_ kv.Iterator[any, any] = FilterIter[any, any]{}
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f FilterIter[K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(f.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (f FilterIter[K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(f.Next, traker)
}

// Next returns the next key/value pair.
// The ok result indicates whether the pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f FilterIter[K, V]) Next() (key K, value V, ok bool) {
	if !(f.next == nil || f.filter == nil) {
		key, value, ok = nextFiltered(f.next, f.filter)
	}
	return key, value, ok
}

// Start is used with for loop construct like 'for i, k, v, ok := i.Start(); ok; k, v, ok = i.Next() { }'
func (f FilterIter[K, V]) Start() (FilterIter[K, V], K, V, bool) {
	return startKvIt[K, V](f)
}

func nextFiltered[K any, V any](next func() (K, V, bool), filter func(K, V) bool) (key K, val V, filtered bool) {
	return First(next, filter)
}
