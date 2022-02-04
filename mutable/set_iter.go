package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

func NewSetIter[k comparable](uniques map[k]struct{}, del func(element k) bool) *SetIter[k] {
	return &SetIter[k]{it.NewKey(uniques), del}
}

type SetIter[k comparable] struct {
	*it.Key[k, struct{}]
	del func(element k) bool
}

var _ c.Iterator[int] = (*SetIter[int])(nil)
var _ Iterator[int] = (*SetIter[int])(nil)

func (iter *SetIter[k]) Next() (k) {
	key, _ := iter.KV.Next()
	return key
}

func (iter *SetIter[k]) Delete() bool {
	return iter.del(iter.Next())
}
