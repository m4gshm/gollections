package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

//NewHead creates the Iter starting at the first element of a sclie.
func NewHead[T any](elements *[]T, del func(int) bool) *Iter[T] {
	return &Iter[T]{elements: elements, current: it.NoStarted, del: del}
}

//NewTail creates the Iter starting at the last element of a sclie.
func NewTail[T any](elements *[]T, del func(int) bool) *Iter[T] {
	return &Iter[T]{elements: elements, current: len(*elements), del: del}
}

//Iter is the Iterator implementation for mutable containers.
type Iter[T any] struct {
	elements *[]T
	current, step  int
	del      func(index int) bool
}

var (
	_ c.Iterator[any]     = (*Iter[any])(nil)
	_ c.PrevIterator[any] = (*Iter[any])(nil)
	_ c.DelIterator[any]  = (*Iter[any])(nil)
)

func (i *Iter[T]) HasNext() bool {
	return it.HasNext(*i.elements, i.current)
}

func (i *Iter[T]) HasPrev() bool {
	return it.HasPrev(*i.elements, i.current)
}

func (i *Iter[T]) GetNext() (T, bool) {
	if i.HasNext() {
		i.current++
		i.step = 1
		return it.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

func (i *Iter[T]) GetPrev() (T, bool) {
	if i.HasPrev() {
		i.current--
		i.step = 0
		return it.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

func (i *Iter[T]) Delete() bool {
	if deleted := i.del(i.current); deleted {
		i.current -= i.step
		// i.deleted = true
		return true
	}
	return false
}

func (i *Iter[T]) DeleteNext() bool {
	if deleted := i.del(i.current + 1); deleted {
		return true
	}
	return false
}

func (i *Iter[T]) DeletePrev() bool {
	if deleted := i.del(i.current - 1); deleted {
		i.current--
		return true
	}
	return false
}