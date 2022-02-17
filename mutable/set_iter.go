package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

//NewSetIter is the default SetIter constructor.
func NewSetIter[K comparable](uniques map[K]struct{}, del func(element K) bool) *SetIter[K] {
	return &SetIter[K]{it.NewKey(uniques), del}
}

//SetIter is the Set Iterator implementation.
type SetIter[K comparable] struct {
	*it.Key[K, struct{}]
	del func(element K) bool
}

var (
	_ c.Iterator[int]    = (*SetIter[int])(nil)
	_ c.DelIterator[int] = (*SetIter[int])(nil)
)

func (i *SetIter[K]) Next() K {
	key, _ := i.KV.Next()
	return key
}

func (i *SetIter[K]) DeleteNext() bool {
	return i.del(i.Next())
}
