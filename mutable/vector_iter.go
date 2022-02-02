package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

func NewIter[T any](elements **[]T, changeMark *int32, del func(int) (bool, error)) *Iter[T] {
	return &Iter[T]{elements: elements, current: it.NoStarted, changeMark: changeMark, del: del}
}

type Iter[T any] struct {
	elements   **[]T
	current    int
	changeMark *int32
	del        func(index int) (bool, error)
}

var _ c.Iterator[any] = (*Iter[any])(nil)
var _ Iterator[any] = (*Iter[any])(nil)

func (i *Iter[T]) HasNext() bool {
	return it.HasNext(*i.elements, &i.current)
}

func (i *Iter[T]) Next() T {
	return it.Get(*i.elements, i.current)
}

func (i *Iter[T]) Delete() (bool, error) {
	pos := i.current
	if deleted, err := i.del(pos); err != nil {
		return false, err
	} else if deleted {
		i.current = pos - 1
		return true, nil
	}
	return false, nil
}
