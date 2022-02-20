package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

//NewSetIter is the default SetIter constructor.
func NewSetIter[K comparable](uniques map[K]struct{}, del func(element K) bool) *SetIter[K] {
	return &SetIter[K]{Key: it.NewKey(uniques), del: del}
}

//SetIter is the Set Iterator implementation.
type SetIter[K comparable] struct {
	*it.Key[K, struct{}]
	del        func(element K) bool
	currentKey K
	ok         bool
}

var (
	_ c.Iterator[int]    = (*SetIter[int])(nil)
	_ c.DelIterator[int] = (*SetIter[int])(nil)
)

func (i *SetIter[K]) Next() (K, bool) {
	key, _, ok := i.KV.Next()
	i.currentKey = key
	i.ok = ok
	return key, ok
}

func (i *SetIter[K]) Delete() bool {
	if i.ok {
		return i.del(i.currentKey)
	}
	return false
}
