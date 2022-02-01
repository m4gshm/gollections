package mutable

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/typ"
)

func NewSetIter[k comparable](uniques map[k]struct{}, del func(element k) (bool, error)) *SetIter[k] {
	return &SetIter[k]{it.NewKey(uniques), del}
}

type SetIter[k comparable] struct {
	*it.Key[k, struct{}]
	del func(element k) (bool, error)
}

var _ typ.Iterator[any] = (*SetIter[any])(nil)
var _ Iterator[any] = (*SetIter[any])(nil)

func (iter *SetIter[k]) Get() (k, error) {
	key, _, err := iter.KV.Get()
	return key, err
}

func (iter *SetIter[k]) Next() k {
	return it.Next[k](iter)
}

func (iter *SetIter[k]) Delete() (bool, error) {
	e, err := iter.Get()
	if err != nil {
		return false, err
	}
	return iter.del(e)
}
