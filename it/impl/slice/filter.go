package slice

import (
	"github.com/m4gshm/gollections/c"
)

//go:notinheap
type Fit[T any] struct {
	elements []T
	by       c.Predicate[T]
	i        int
}

var _ c.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) Next() (T, bool) {
	return nextArrayElem(s.elements, s.by, &s.i)
}

func (s *Fit[T]) Cap() int {
	return len(s.elements)
}

func (s Fit[T]) R() *Fit[T] {
	return &s
	// return notsafe.Noescape(&s)
}
