package slice

import (
	"github.com/m4gshm/gollections/c"
)

type Fit[T any] struct {
	Elements []T
	By       c.Predicate[T]
	i       int
}

var _ c.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) GetNext() (T, bool) {
	return nextArrayElem(s.Elements, s.By, &s.i)
}
