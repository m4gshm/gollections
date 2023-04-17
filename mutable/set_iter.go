package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// NewSetIter creates SetIter instance.
func NewSetIter[K comparable](uniques map[K]struct{}, del func(element K)) *SetIter[K] {
	return &SetIter[K]{Key: iter.NewKey(uniques), del: del}
}

// SetIter is the Set Iterator implementation.
type SetIter[K comparable] struct {
	iter.Key[K, struct{}]
	del        func(element K)
	currentKey K
	ok         bool
}

var (
	_ c.Iterator[int]    = (*SetIter[int])(nil)
	_ c.DelIterator[int] = (*SetIter[int])(nil)
)

func (i *SetIter[K]) Next() (K, bool) {
	key, _, ok := i.EmbedMapKVIter.Next()
	i.currentKey = key
	i.ok = ok
	return key, ok
}

func (i *SetIter[K]) Delete() {
	if i.ok {
		i.del(i.currentKey)
	}
}
