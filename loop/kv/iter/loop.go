package iter

import "github.com/m4gshm/gollections/c"

// New creates an LoopKVIter instance that loops over key\value elements of a source.
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the one.
func New[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) LoopKVIter[S, K, V] {
	return LoopKVIter[S, K, V]{source: source, hasNext: hasNext, getNext: getNext}
}

// LoopKVIter - universal key\value iterator implementation
type LoopKVIter[S, K, V any] struct {
	source  S
	hasNext func(S) bool
	getNext func(S) (K, V, error)
	abort   error
}

var (
	_ c.KVIterator[any, any]          = (*LoopKVIter[any, any, any])(nil)
	_ c.KVIteratorBreakable[any, any] = (*LoopKVIter[any, any, any])(nil)
)

// Next implements c.KVIterator
func (i *LoopKVIter[S, K, V]) Next() (K, V, bool) {
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
func (i *LoopKVIter[S, K, V]) Error() error {
	return i.abort
}
