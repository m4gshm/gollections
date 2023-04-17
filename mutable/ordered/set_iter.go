package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// NewSetIter creates a set's iterator.
func NewSetIter[T any](elements *[]T, del func(v T)) *SetIter[T] {
	return &SetIter[T]{elements: elements, current: iter.NoStarted, del: del}
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

func (i *SetIter[T]) Next() (T, bool) {
	if iter.HasNext(*i.elements, i.current) {
		i.current++
		return iter.Gett(*i.elements, i.current)
	}
	var no T
	return no, false
}

func (i *SetIter[T]) Cap() int {
	return len(*i.elements)
}

func (i *SetIter[T]) Delete() {
	if v, ok := iter.Gett(*i.elements, i.current); ok {
		i.current--
		i.del(v)
	}
}
