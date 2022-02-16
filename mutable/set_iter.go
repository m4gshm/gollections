package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

//NewSetIter is the default SetIter constructor.
func NewSetIter[k comparable](uniques map[k]struct{}, del func(element k) bool) *SetIter[k] {
	return &SetIter[k]{it.NewKey(uniques), del}
}

//SetIter is the Set Iterator implementation.
type SetIter[k comparable] struct {
	*it.Key[k, struct{}]
	del func(element k) bool
}

var (
	_ c.Iterator[int]    = (*SetIter[int])(nil)
	_ c.DelIterator[int] = (*SetIter[int])(nil)
)

func (i *SetIter[k]) Get() k {
	key, _ := i.KV.Get()
	return key
}

func (i *SetIter[k]) Delete() bool {
	return i.del(i.Get())
}
