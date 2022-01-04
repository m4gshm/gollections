package set

import (
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/mutable/iter"
	"github.com/m4gshm/container/typ"
)

type RefIter[T any] struct {
	*iter.Deleteable[*T]
}

var _ typ.Iterator[any] = (*RefIter[any])(nil)
var _ mutable.Iterator[any] = (*RefIter[any])(nil)

func (i *RefIter[T]) HasNext() bool {
	return i.Deleteable.HasNext()

}
func (i *RefIter[T]) Get() T {
	return *i.Deleteable.Get()
}

func (i *RefIter[T]) Delete() bool {
	return i.Deleteable.Delete()
}
