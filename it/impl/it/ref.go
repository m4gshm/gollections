package it

import (
	"github.com/m4gshm/gollections/typ"
)

func NewRef[T any](elements []*T) *RefIter[T] {
	return &RefIter[T]{New(elements)}
}

type RefIter[T any] struct {
	typ.Iterator[*T]
}

var _ typ.Iterator[any] = (*RefIter[any])(nil)

func (i *RefIter[T]) Get() (T, error) {
	v, err := i.Iterator.Get()
	if err != nil {
		var no T
		return no, err
	}
	return *v, nil
}
