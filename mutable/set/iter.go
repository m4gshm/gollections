package set

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/typ"
)

func NewIter[k comparable](uniques map[k]struct{}, del func(element k) (bool, error)) *Iter[k] {
	return &Iter[k]{it.NewKey(uniques), del}
}

type Iter[k comparable] struct {
	*it.Key[k, struct{}]
	del func(element k) (bool, error)
}

var _ typ.Iterator[any] = (*Iter[any])(nil)
var _ mutable.Iterator[any] = (*Iter[any])(nil)

func (iter *Iter[k]) Get() (k, error) {
	key, _, err := iter.KV.Get()
	return key, err
}

func (iter *Iter[k]) Next() k {
	return it.Next[k](iter)
}

func (iter *Iter[k]) Delete() (bool, error) {
	e, err := iter.Get()
	if err != nil {
		return false, err
	}
	return iter.del(e)
}
