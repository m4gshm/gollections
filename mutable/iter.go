package mutable

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)


type OrderIter[T any] struct {
	*iter.Deleteable[*T]
}

var _ typ.Iterator[any] = (*OrderIter[any])(nil)

func (i *OrderIter[T]) HasNext() bool {
	return i.Iter.HasNext()

}
func (i *OrderIter[T]) Get() T {
	return *i.Iter.Get()
}

func (i *OrderIter[T]) Delete() bool {
	return i.Deleteable.Delete()
}