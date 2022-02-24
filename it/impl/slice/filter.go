package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/notsafe"
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

//Experimental
//must be inlined
//DON'T USE IN PROD
func (s Fit[T]) R() *Fit[T] {
	return notsafe.Noescape(&s)
}
