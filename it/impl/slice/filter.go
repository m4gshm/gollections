package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

type Fit[T any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	by       c.Predicate[T]
}

var _ c.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) Next() (T, bool) {
	return nextFiltered(s.array, s.size, s.elemSize, s.by, &s.i)
}

func (s *Fit[T]) Cap() int {
	return s.size
}
