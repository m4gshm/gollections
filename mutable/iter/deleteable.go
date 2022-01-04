package iter

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
)

const NoStarted = -1

func NewDeleteable[T any](elements *[]T, changeMark *int32, del func(v T) bool) *Deleteable[T] {
	return &Deleteable[T]{elements: elements, current: NoStarted, changeMark: changeMark, del: del}
}

type Deleteable[T any] struct {
	elements   *[]T
	err        error
	current    int
	changeMark *int32
	del        func(v T) bool
}

var _ typ.Iterator[any] = (*Deleteable[any])(nil)
var _ mutable.Iterator[any] = (*Deleteable[any])(nil)

func (i *Deleteable[T]) HasNext() bool {
	return iter.HasNext(*i.elements, &i.current, &i.err)
}
func (i *Deleteable[T]) Get() T {
	return iter.Get(i.current, *i.elements, i.err)
}

func (i *Deleteable[T]) Delete() bool {
	pos := i.current
	if deleted := i.del(i.Get()); deleted {
		i.current = pos - 1
		return true
	}
	return false
}

func (s *Deleteable[T]) Err() error {
	return s.err
}

