package iter

import (
	"github.com/m4gshm/gollections/c"
)

// ConvertFitIter is the Converter with elements filtering.
type ConvertFitIter[From, To, IT any] struct {
	iterator IT
	next     func() (From, bool)
	by       func(From) To
	filter   func(From) bool
}

var (
	_ c.Iterator[any] = (*ConvertFitIter[any, any, any])(nil)
)

func (c *ConvertFitIter[From, To, IT]) Next() (t To, ok bool) {
	if c == nil || c.by == nil {
		return
	}
	if V, ok := nextFiltered(c.next, c.filter); ok {
		return c.by(V), true
	}
	return
}

// ConvertIter is the iterator wrapper implementation applying a converter to all iterable elements.
type ConvertIter[From, To any, IT any] struct {
	iterator IT
	next     func() (From, bool)
	by       func(From) To
}

var (
	_ c.Iterator[any] = (*ConvertIter[any, any, any])(nil)
)

func (c *ConvertIter[From, To, IT]) Next() (t To, ok bool) {
	if c == nil || c.by == nil {
		return
	}
	if v, ok := c.next(); ok {
		return c.by(v), true
	}
	return
}

// ConvertKVIter is the iterator wrapper implementation applying a converter to all iterable key/value elements.
type ConvertKVIter[K, V any, IT c.KVIterator[K, V], K2, V2 any, C func(K, V) (K2, V2)] struct {
	iterator IT
	by       C
}

var (
	_ c.KVIterator[any, any] = (*ConvertKVIter[any, any, c.KVIterator[any, any], any, any, func(any, any) (any, any)])(nil)
)

func (c *ConvertKVIter[K, V, IT, K2, V2, C]) Next() (k K2, v V2, ok bool) {
	if c == nil || c.by == nil {
		return
	}
	if K, V, ok := c.iterator.Next(); ok {
		k2, v2 := c.by(K, V)
		return k2, v2, true
	}
	return
}
