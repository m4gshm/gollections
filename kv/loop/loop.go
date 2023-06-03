package loop

import (
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/loop"
)

// NewIter creates an Iter instance that loops over key\value elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the one.
func NewIter[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) Iter[S, K, V] {
	return Iter[S, K, V]{source: source, hasNext: hasNext, getNext: getNext}
}

// Iter - universal key\value iterator implementation
type Iter[S, K, V any] struct {
	source  S
	hasNext func(S) bool
	getNext func(S) (K, V, error)
	abort   error
}

var (
	_ kv.Iterator[any, any] = (*Iter[any, any, any])(nil)
)

// Track takes key, value pairs retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *Iter[S, K, V]) Track(traker func(key K, value V) error) error {
	return loop.Track(i.Next, traker)
}

// TrackEach takes all key, value pairs retrieved by the iterator
func (i *Iter[S, K, V]) TrackEach(traker func(key K, value V)) {
	loop.TrackEach(i.Next, traker)
}

// Next implements kv.KVIterator
func (i *Iter[S, K, V]) Next() (K, V, bool) {
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

// Error implements kv.KVIteratorBreakable
func (i *Iter[S, K, V]) Error() error {
	return i.abort
}

func (i *Iter[S, K, V]) Start() (*Iter[S, K, V], K, V, bool) {
	return startKvIt[K, V](i)
}
