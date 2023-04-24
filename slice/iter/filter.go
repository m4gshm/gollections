package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
)

// FitIter is the array based Iterator implementation that provides filtering of elements by a Predicate.
type FitIter[T any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	filter   func(T) bool
}

var _ c.Iterator[any] = (*FitIter[any])(nil)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (s *FitIter[T]) Next() (T, bool) {
	return nextFiltered(s.array, s.size, s.elemSize, s.filter, &s.i)
}

// Cap returns the iterator capacity
func (s *FitIter[T]) Cap() int {
	return s.size
}
