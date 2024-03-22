package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// NewHead instantiates Iter starting at the first element of a slice.
func NewHead[TS ~[]T, T any](elements *TS, del func(int) bool) SliceIter[T] {
	return SliceIter[T]{elements: slice.UpcastRef(elements), current: slice.IterNoStarted, del: del}
}

// NewTail instantiates Iter starting at the last element of a slice.
func NewTail[TS ~[]T, T any](elements *TS, del func(int) bool) SliceIter[T] {
	if elements == nil {
		return SliceIter[T]{}
	}
	return SliceIter[T]{elements: slice.UpcastRef(elements), current: len(*elements), del: del}
}

// SliceIter is the Iterator implementation for mutable containers.
type SliceIter[T any] struct {
	elements      *[]T
	current, step int
	del           func(index int) bool
}

var (
	_ c.Iterator[any]     = (*SliceIter[any])(nil)
	_ c.PrevIterator[any] = (*SliceIter[any])(nil)
	_ c.DelIterator[any]  = (*SliceIter[any])(nil)
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (i *SliceIter[T]) For(walker func(element T) error) error {
	return loop.For(i.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (i *SliceIter[T]) ForEach(walker func(element T)) {
	loop.ForEach(i.Next, walker)
}

// HasNext checks the next element existing
func (i *SliceIter[T]) HasNext() bool {
	if i == nil || i.elements == nil {
		return false
	}
	return slice.HasNext(*i.elements, i.current)
}

// HasPrev checks the previous element existing
func (i *SliceIter[T]) HasPrev() bool {
	if i == nil || i.elements == nil {
		return false
	}
	return slice.HasPrev(*i.elements, i.current)
}

// GetNext returns the next element
func (i *SliceIter[T]) GetNext() (t T) {
	if i != nil {
		t, _ = i.Next()
	}
	return
}

// GetPrev returns the previous element
func (i *SliceIter[T]) GetPrev() (t T) {
	if i != nil {
		t, _ = i.Prev()
	}
	return
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *SliceIter[T]) Next() (T, bool) {
	if i.HasNext() {
		i.current++
		i.step = 1
		return slice.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

// Prev returns the previous element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *SliceIter[T]) Prev() (T, bool) {
	if i.HasPrev() {
		i.current--
		i.step = 0
		return slice.Get(*i.elements, i.current), true
	}
	var no T
	return no, false
}

// Get returns the current element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *SliceIter[T]) Get() (t T, ok bool) {
	if i == nil || i.elements == nil {
		return t, ok
	}
	current := i.current
	elements := *i.elements
	if slice.IsValidIndex(len(elements), current) {
		return elements[current], true
	}
	return t, ok
}

// Size returns the iterator capacity
func (i *SliceIter[T]) Size() int {
	if i == nil || i.elements == nil {
		return 0
	}
	return len(*i.elements)
}

// Delete deletes the current element
func (i *SliceIter[T]) Delete() {
	if i == nil {
		return
	} else if deleted := i.del(i.current); deleted {
		i.current -= i.step
	}
}

// DeleteNext deletes the next element if it exists
func (i *SliceIter[T]) DeleteNext() bool {
	if i == nil {
		return false
	}
	return i.del(i.current + 1)
}

// DeletePrev deletes the previos element if it exists
func (i *SliceIter[T]) DeletePrev() bool {
	if i == nil {
		return false
	} else if deleted := i.del(i.current - 1); deleted {
		i.current--
		return true
	}
	return false
}
