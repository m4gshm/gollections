package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// NewHead instantiates Iter starting at the first element of a slice.
func NewHead[TS ~[]T, T any](elements *TS, del func(int) bool) Iter[TS, T] {
	return Iter[TS, T]{elements: elements, current: iter.NoStarted, del: del}
}

// NewTail instantiates Iter starting at the last element of a slice.
func NewTail[TS ~[]T, T any](elements *TS, del func(int) bool) Iter[TS, T] {
	if elements == nil {
		return Iter[TS, T]{}
	}
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
	if i == nil || i.elements == nil {
		return false
	}
	return iter.HasNext(*i.elements, i.current)
}

func (i *Iter[TS, T]) HasPrev() bool {
	if i == nil || i.elements == nil {
		return false
	}
	return iter.HasPrev(*i.elements, i.current)
}

func (i *Iter[TS, T]) GetNext() (t T) {
	if i != nil {
		t, _ = i.Next()
	}
	return
}

func (i *Iter[TS, T]) GetPrev() (t T) {
	if i != nil {
		t, _ = i.Prev()
	}
	return
}

func (i *Iter[TS, T]) Next() (T, bool) {
	if i.HasNext() {
		i.current++
		i.step = 1
		return iter.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

func (i *Iter[TS, T]) Prev() (T, bool) {
	if i.HasPrev() {
		i.current--
		i.step = 0
		return iter.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

func (i *Iter[TS, T]) Get() (t T, ok bool) {
	if i == nil || i.elements == nil {
		return t, ok
	}
	current := i.current
	elements := *i.elements
	if iter.IsValidIndex(len(elements), current) {
		return elements[current], true
	}
	return t, ok
}

func (i *Iter[TS, T]) Cap() int {
	if i == nil || i.elements == nil {
		return 0
	}
	return len(*i.elements)
}

func (i *Iter[TS, T]) Delete() {
	if i == nil {
		return
	} else if deleted := i.del(i.current); deleted {
		i.current -= i.step
	}
}

func (i *Iter[TS, T]) DeleteNext() bool {
	if i == nil {
		return false
	}
	return i.del(i.current + 1)
}

func (i *Iter[TS, T]) DeletePrev() bool {
	if i == nil {
		return false
	} else if deleted := i.del(i.current - 1); deleted {
		i.current--
		return true
	}
	return false
}
