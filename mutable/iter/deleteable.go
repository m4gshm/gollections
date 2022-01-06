package it

import (
	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
)

const NoStarted = -1

func NewDeleteable[T any](elements *[]T, changeMark *int32, del func(v T) (bool, error)) *Deleteable[T] {
	return &Deleteable[T]{elements: elements, current: NoStarted, changeMark: changeMark, del: del}
}

type Deleteable[T any] struct {
	elements   *[]T
	err        error
	current    int
	changeMark *int32
	del        func(v T) (bool, error)
}

var _ typ.Iterator[any] = (*Deleteable[any])(nil)
var _ mutable.Iterator[any] = (*Deleteable[any])(nil)

func (i *Deleteable[T]) HasNext() bool {
	return it.HasNext(*i.elements, &i.current, &i.err)
}
func (i *Deleteable[T]) Get() T {
	return it.Get(i.current, *i.elements, i.err)
}

func (i *Deleteable[T]) Delete() (bool, error) {
	pos := i.current
	if deleted, err := i.del(i.Get()); err != nil {
		return false, err
	} else if deleted {
		i.current = pos - 1
		return true, nil
	}
	return false, nil
}

func (s *Deleteable[T]) Err() error {
	return s.err
}
