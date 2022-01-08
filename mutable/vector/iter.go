package vector

import (
	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
)

func NewIter[T any](elements **[]T, changeMark *int32, del func(int) (bool, error)) *Iter[T] {
	return &Iter[T]{elements: elements, current: it.NoStarted, changeMark: changeMark, del: del}
}

type Iter[T any] struct {
	elements   **[]T
	err        error
	current    int
	changeMark *int32
	del        func(index int) (bool, error)
}

var _ typ.Iterator[any] = (*Iter[any])(nil)
var _ mutable.Iterator[any] = (*Iter[any])(nil)

func (i *Iter[T]) HasNext() bool {
	return it.HasNext(**i.elements, &i.current, &i.err)
}
func (i *Iter[T]) Get() T {
	return it.Get(i.current, **i.elements, i.err)
}

func (i *Iter[T]) Delete() (bool, error) {
	pos := i.current
	if deleted, err := i.del(pos); err != nil {
		return false, err
	} else if deleted {
		i.current = pos - 1
		return true, nil
	}
	return false, nil
}

func (s *Iter[T]) Err() error {
	return s.err
}

type RefIter[T any] struct {
	*Iter[*T]
}

var _ typ.Iterator[any] = (*RefIter[any])(nil)
var _ mutable.Iterator[any] = (*RefIter[any])(nil)

func (i *RefIter[T]) HasNext() bool {
	return i.Iter.HasNext()
}

func (i *RefIter[T]) Get() T {
	return *i.Iter.Get()
}

func (i *RefIter[T]) Delete() (bool, error) {
	return i.Iter.Delete()
}
