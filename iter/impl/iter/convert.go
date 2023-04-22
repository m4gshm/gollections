package iter

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertFitIter is the Converter with elements filtering.
type ConvertFitIter[From, To any] struct {
	next   func() (From, bool)
	by     func(From) To
	filter func(From) bool
}

var (
	_ c.Iterator[any] = (*ConvertFitIter[any, any])(nil)
	_ c.Iterator[any] = ConvertFitIter[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertFitIter[From, To]) Next() (t To, ok bool) {
	if next, filter := c.next, c.filter; next != nil && filter != nil {
		if f, ok := nextFiltered(next, filter); ok {
			return c.by(f), true
		}
	}
	return t, false
}

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable elements.
type ConvertIter[From, To any] struct {
	next      func() (From, bool)
	converter func(From) To
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any])(nil)
	_ c.Iterator[any] = ConvertIter[any, any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (c ConvertIter[From, To]) Next() (t To, ok bool) {
	if next := c.next; next != nil {
		if v, ok := next(); ok {
			return c.converter(v), true
		}
	}
	return t, false
}

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
