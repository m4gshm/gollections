package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
)

// NewHead instantiates Iter starting at the first element of a slice.
func NewHead[TS ~[]T, T any](elements *TS, del func(int) bool) Iter[TS, T] {
	return Iter[TS, T]{elements: elements, current: it.NoStarted, del: del}
}

// NewTail instantiates Iter starting at the last element of a slice.
func NewTail[TS ~[]T, T any](elements *TS, del func(int) bool) Iter[TS, T] {
	return Iter[TS, T]{elements: elements, current: len(*elements), del: del}
}

// Iter is the Iterator implementation for mutable containers.
type Iter[TS ~[]T, T any] struct {
	elements      *TS
	current, step int
	del           func(index int) bool
}

var (
	_ c.Iterator[any]     = (*Iter[[]any, any])(nil)
	_ c.PrevIterator[any] = (*Iter[[]any, any])(nil)
	_ c.DelIterator[any]  = (*Iter[[]any, any])(nil)
)

func (i *Iter[TS, T]) HasNext() bool {
	return it.HasNext(*i.elements, i.current)
}

func (i *Iter[TS, T]) HasPrev() bool {
	return it.HasPrev(*i.elements, i.current)
}

func (i *Iter[TS, T]) GetNext() T {
	t, _ := i.Next()
	return t
}

func (i *Iter[TS, T]) GetPrev() T {
	t, _ := i.Prev()
	return t
}

func (i *Iter[TS, T]) Next() (T, bool) {
	if i.HasNext() {
		i.current++
		i.step = 1
		return it.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

func (i *Iter[TS, T]) Prev() (T, bool) {
	if i.HasPrev() {
		i.current--
		i.step = 0
		return it.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

func (i *Iter[TS, T]) Get() (T, bool) {
	current := i.current
	elements := *i.elements
	if it.IsValidIndex(len(elements), current) {
		return elements[current], true
	}
	var no T
	return no, false
}

func (i *Iter[TS, T]) Cap() int {
	return len(*i.elements)
}

func (i *Iter[TS, T]) Delete() bool {
	if deleted := i.del(i.current); deleted {
		i.current -= i.step
		// i.deleted = true
		return true
	}
	return false
}

func (i *Iter[TS, T]) DeleteNext() bool {
	if deleted := i.del(i.current + 1); deleted {
		return true
	}
	return false
}

func (i *Iter[TS, T]) DeletePrev() bool {
	if deleted := i.del(i.current - 1); deleted {
		i.current--
		return true
	}
	return false
}
