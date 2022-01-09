package set

import (
	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
)

func NewIter[T any](elements *[]*T, changeMark *int32, del func(v T) (bool, error)) *Iter[T] {
	return &Iter[T]{elements: elements, current: it.NoStarted, changeMark: changeMark, del: del}
}

type Iter[T any] struct {
	elements   *[]*T
	err        error
	current    int
	changeMark *int32
	del        func(v T) (bool, error)
}

var _ typ.Iterator[any] = (*Iter[any])(nil)
var _ mutable.Iterator[any] = (*Iter[any])(nil)

func (i *Iter[T]) HasNext() bool {
	return it.HasNext(*i.elements, &i.current, &i.err)
}

func (i *Iter[T]) Next() (T, error) {
	v, err := it.Next(i.current, *i.elements, i.err)
	if err != nil {
		var no T
		return no, err
	}
	return *v, nil
}

func (i *Iter[T]) Get() T {
	return *it.Get(i.current, *i.elements, i.err)
}

func (i *Iter[T]) Delete() (bool, error) {
	pos := i.current
	if e, err := i.Next(); err != nil {
		return false, err
	} else if deleted, err := i.del(e); err != nil {
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
