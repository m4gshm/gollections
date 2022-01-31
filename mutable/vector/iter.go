package vector

import (
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/typ"
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
	return it.HasNext(*i.elements, &i.current, &i.err)
}

func (i *Iter[T]) Get() (T, error) {
	return it.Get(*i.elements, i.current, i.err)
}

func (s *Iter[T]) Next() T {
	return it.Next[T](s)
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

type RefIter[T any] struct {
	*Iter[*T]
}

var _ typ.Iterator[any] = (*RefIter[any])(nil)
var _ mutable.Iterator[any] = (*RefIter[any])(nil)

func (i *RefIter[T]) HasNext() bool {
	return i.Iter.HasNext()
}

func (i *RefIter[T]) Get() (T, error) {
	v, err := i.Iter.Get()
	if err != nil {
		var no T
		return no, err
	}
	return *v, nil
}

func (s *RefIter[T]) Next() T {
	return it.Next[T](s)
}

func (i *RefIter[T]) Delete() (bool, error) {
	return i.Iter.Delete()
}
