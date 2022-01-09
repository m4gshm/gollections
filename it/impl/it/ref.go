package it

import (
	"github.com/m4gshm/container/typ"
)

func NewRef[T any](elements []*T) *RefIter[T] {
	return &RefIter[T]{New(elements)}
}

type RefIter[T any] struct {
	typ.Iterator[*T]
}

var _ typ.Iterator[any] = (*RefIter[any])(nil)

func (i *RefIter[T]) Get() T {
	return *i.Iterator.Get()
}

func (i *RefIter[T]) Next() (T, error) {
	v, err := i.Iterator.Next()
	if err != nil {
		var no T
		return no, err
	}
	return *v, nil
}
