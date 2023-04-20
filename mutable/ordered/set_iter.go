package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// NewSetIter creates a set's iterator.
func NewSetIter[T any](elements *[]T, del func(v T)) SetIter[T] {
	return SetIter[T]{elements: elements, current: iter.NoStarted, del: del}
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

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (i *SetIter[T]) Next() (t T, ok bool) {
	if !(i == nil || i.elements == nil) {
		if iter.HasNext(*i.elements, i.current) {
			i.current++
			return iter.Gett(*i.elements, i.current)
		}
	}
	return t, ok
}

// Cap returns the iterator capacity
func (i *SetIter[T]) Cap() (capacity int) {
	if !(i == nil || i.elements == nil) {
		capacity = len(*i.elements)
	}
	return capacity
}

// Delete deletes the current element
func (i *SetIter[T]) Delete() {
	if !(i == nil || i.elements == nil) {
		if v, ok := iter.Gett(*i.elements, i.current); ok {
			i.current--
			i.del(v)
		}
	}
}
