package it

import (
	"github.com/m4gshm/gollections/typ"
)

func NewRef[T any](elements *[]*T) *RefIter[T] {
	return WrapRef(NewP(elements))
}

func WrapRef[T any](it *PIter[*T]) *RefIter[T] {
	return &RefIter[T]{it}
}

type RefIter[T any] struct {
	// typ.Iterator[*T]
	*PIter[*T]
}

var _ typ.Iterator[any] = (*RefIter[any])(nil)

func (i *RefIter[T]) Get() (T, error) {
	v, err := i.PIter.Get()
	if err != nil {
		var no T
		return no, err
	}
	return *v, nil
}

func (i *RefIter[T]) Next() T {
	return Next[T](i)
}
