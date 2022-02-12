package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

func NewIter[T any](elements *[]T, del func(int) bool) *Iter[T] {
	return &Iter[T]{elements: elements, current: it.NoStarted, del: del}
}

type Iter[T any] struct {
	elements *[]T
	current  int
	del      func(index int) bool
}

var _ c.Iterator[any] = (*Iter[any])(nil)
var _ Iterator[any] = (*Iter[any])(nil)

func (i *Iter[T]) HasNext() bool {
	if it.HasNext(*i.elements, i.current) {
		i.current++
		return true
	}
	return false
}

func (i *Iter[T]) Get() T {
	return it.Get(*i.elements, i.current)
}

func (i *Iter[T]) Delete() bool {
	pos := i.current
	if deleted := i.del(pos); deleted {
		i.current = pos - 1
		return true
	}
	return false
}
