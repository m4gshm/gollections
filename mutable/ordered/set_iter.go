package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
)

func NewSetIter[T any](elements *[]T, del func(v T) bool) *SetIter[T] {
	return &SetIter[T]{elements: elements, current: it.NoStarted, del: del}
}

type SetIter[T any] struct {
	elements *[]T
	current  int
	del      func(v T) bool
}

var _ c.Iterator[any] = (*SetIter[any])(nil)
var _ mutable.Iterator[any] = (*SetIter[any])(nil)

func (i *SetIter[T]) HasNext() bool {
	if it.HasNext(*i.elements, i.current) {
		i.current++
		return true
	}
	return false
}

func (i *SetIter[T]) Get() T {
	return it.Get(*i.elements, i.current)
}

func (i *SetIter[T]) Delete() bool {
	pos := i.current
	if deleted := i.del(i.Get()); deleted {
		i.current = pos - 1
		return true
	}
	return false
}
