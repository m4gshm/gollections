package loop

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// NewIter creates an KVIter instance that loops over key\value elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the one.
func NewIter[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) KVIter[S, K, V] {
	return KVIter[S, K, V]{source: source, hasNext: hasNext, getNext: getNext}
}

// KVIter - universal key\value iterator implementation
type KVIter[S, K, V any] struct {
	source  S
	hasNext func(S) bool
	getNext func(S) (K, V, error)
	abort   error
}

var (
	_ c.KVIterator[any, any]          = (*KVIter[any, any, any])(nil)
	_ c.KVIteratorBreakable[any, any] = (*KVIter[any, any, any])(nil)
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (kv *KVIter[S, K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(kv.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (kv *KVIter[S, K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(kv.Next, traker)
}

// Next implements c.KVIterator
func (i *KVIter[S, K, V]) Next() (K, V, bool) {
	if i.abort == nil && i.hasNext(i.source) {
		k, v, err := i.getNext(i.source)
		if err == nil {
			return k, v, true
		}
		i.abort = err
	}
	var k K
	var v V
	return k, v, false
}

// Error implements c.KVIteratorBreakable
func (i *KVIter[S, K, V]) Error() error {
	return i.abort
}
