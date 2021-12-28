package immutable

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

type OrderIter[T any] struct {
	*iter.Iter[*T]
}

var _ typ.Iterator[interface{}] = (*OrderIter[interface{}])(nil)

func (i *OrderIter[T]) HasNext() bool {
	return i.Iter.HasNext()
}

func (i *OrderIter[T]) Get() T {
	return *i.Iter.Get()
}
