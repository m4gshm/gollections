package set

import (
	"github.com/m4gshm/container/typ"
)

type OrderIter[T any] struct {
	 typ.Iterator[*T]
}

var _ typ.Iterator[any] = (*OrderIter[any])(nil)

func (i *OrderIter[T]) Get() T {
	return *i.Iterator.Get()
}