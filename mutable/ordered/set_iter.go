package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
)

func NewSetIter[T any](elements *[]T, changeMark *int32, del func(v T) (bool, error)) *SetIter[T] {
	return &SetIter[T]{elements: elements, current: it.NoStarted, changeMark: changeMark, del: del}
}

type SetIter[T any] struct {
	elements   *[]T
	err        error
	current    int
	changeMark *int32
	del        func(v T) (bool, error)
}

var _ c.Iterator[any] = (*SetIter[any])(nil)
var _ mutable.Iterator[any] = (*SetIter[any])(nil)

func (i *SetIter[T]) HasNext() bool {
	return it.HasNext(i.elements, &i.current)
}

func (i *SetIter[T]) Next() T {
	return it.Get(i.elements, i.current)
}

func (i *SetIter[T]) Delete() (bool, error) {
	pos := i.current
	if deleted, err := i.del(i.Next()); err != nil {
		return false, err
	} else if deleted {
		i.current = pos - 1
		return true, nil
	}
	return false, nil
}
