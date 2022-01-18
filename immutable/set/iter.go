package set

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/typ"
)

type SetIter[k comparable] struct {
	*it.KV[k, struct{}]
}

var _ typ.Iterator[any] = (*SetIter[any])(nil)

func (iter *SetIter[k]) Get() (k, error) {
	key, _, err := iter.KV.Get()
	return key, err
}
