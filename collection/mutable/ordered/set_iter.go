package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/slice"
)

// NewSetIter creates a set's iterator.
func NewSetIter[T any](elements *[]T, del func(v T)) SetIter[T] {
	return SetIter[T]{elements: elements, current: slice.IterNoStarted, del: del}
}

// SetIter set iterator
type SetIter[T any] struct {
	elements *[]T
	current  int
	del      func(v T)
}

var (
	_ c.Iterator[any]    = (*SetIter[any])(nil)
	_ c.DelIterator[any] = (*SetIter[any])(nil)
)

// All is used to iterate through the collection using `for e := range`.
func (i *SetIter[T]) All(consumer func(element T) bool) {
	loop.All(i.Next, consumer)
}

// For takes elements retrieved by the iterator. Can be interrupt by returning Break
func (i *SetIter[T]) For(consumer func(element T) error) error {
	return loop.For(i.Next, consumer)
}

// ForEach takes all elements retrieved by the iterator.
func (i *SetIter[T]) ForEach(consumer func(element T)) {
	loop.ForEach(i.Next, consumer)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *SetIter[T]) Next() (t T, ok bool) {
	if !(i == nil || i.elements == nil) {
		if slice.HasNext(*i.elements, i.current) {
			i.current++
			return slice.Gett(*i.elements, i.current)
		}
	}
	return t, ok
}

// Size returns the iterator capacity
func (i *SetIter[T]) Size() (capacity int) {
	if !(i == nil || i.elements == nil) {
		capacity = len(*i.elements)
	}
	return capacity
}

// Delete deletes the current element
func (i *SetIter[T]) Delete() {
	if !(i == nil || i.elements == nil) {
		if v, ok := slice.Gett(*i.elements, i.current); ok {
			i.current--
			i.del(v)
		}
	}
}
