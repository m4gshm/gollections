package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

func NewHead[T any](elements *[]T, del func(int) bool) *Iter[T] {
	return &Iter[T]{elements: elements, current: it.NoStarted, del: del}
}

func NewTail[T any](elements *[]T, del func(int) bool) *Iter[T] {
	return &Iter[T]{elements: elements, current: len(*elements), del: del}
}

type Iter[T any] struct {
	current, step int
	elements      *[]T
	del           func(index int) bool
}

var (
	_ c.Iterator[any] = (*Iter[any])(nil)
	_ Iterator[any]   = (*Iter[any])(nil)
)

func (i *Iter[T]) HasNext() bool {
	if it.HasNext(*i.elements, i.current) {
		i.step = 1
		i.current++
		return true
	}
	return false
}

func (i *Iter[T]) HasPrev() bool {
	if it.HasPrev(*i.elements, i.current) {
		i.step = 0
		i.current--
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
		i.current = pos - i.step
		return true
	}
	return false
}
