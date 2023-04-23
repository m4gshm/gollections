package iter

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertKVIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertKVIter[K, V any, K2, V2 any, C func(K, V) (K2, V2)] struct {
	next func() (K, V, bool)
	by   C
}

var (
	_ c.KVIterator[any, any] = (*ConvertKVIter[any, any, any, any, func(any, any) (any, any)])(nil)
	_ c.KVIterator[any, any] = ConvertKVIter[any, any, any, any, func(any, any) (any, any)]{}
)

// Next returns the next key/value pair.
// The ok result indicates whether an pair was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertKVIter[K, V, K2, V2, C]) Next() (k2 K2, v2 V2, ok bool) {
	if next := c.next; next != nil {
		if K, V, ok := next(); ok {
			k2, v2 = c.by(K, V)
			return k2, v2, true
		}
	}
	return k2, v2, false
}
