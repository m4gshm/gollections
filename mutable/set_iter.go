package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/map_/iter"
)

// NewSetIter creates SetIter instance.
func NewSetIter[K comparable](uniques map[K]struct{}, del func(element K)) SetIter[K] {
	return SetIter[K]{KeyIter: iter.NewKey(uniques), del: del}
}

// SetIter is the Set Iterator implementation.
type SetIter[K comparable] struct {
	iter.KeyIter[K, struct{}]
	del        func(element K)
	currentKey K
	ok         bool
}

var (
	_ c.Iterator[int]    = (*SetIter[int])(nil)
	_ c.DelIterator[int] = (*SetIter[int])(nil)
)

// Next returns the next element if it exists
func (i *SetIter[K]) Next() (key K, ok bool) {
	if i != nil {
		key, _, ok = i.MapIter.Next()
		i.currentKey = key
		i.ok = ok
	}
	return key, ok
}

// Delete deletes the current element
func (i *SetIter[K]) Delete() {
	if i != nil && i.ok {
		i.del(i.currentKey)
	}
}
