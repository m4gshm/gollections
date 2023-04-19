package slice

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

// Fit is the array based Iterator implementation that provides filtering of elements by a Predicate.
type Fit[T any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	by       func(T) bool
}

var _ c.Iterator[any] = (*Fit[any])(nil)

func (s *Fit[T]) Next() (T, bool) {
	return nextFiltered(s.array, s.size, s.elemSize, s.by, &s.i)
}

func (s *Fit[T]) Cap() int {
	return s.size
}
