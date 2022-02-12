package slice

import (
	"github.com/m4gshm/gollections/c"
)

type Fit[T any] struct {
	Elements []T
	By       c.Predicate[T]

	current T
	i       int
}

var _ c.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) HasNext() bool {
	v, ok := nextArrayElem(s.Elements, s.By, &s.i)
	if ok {
		s.current = v
	}
	return ok
}

func (s *Fit[T]) Get() T {
	return s.current
}
