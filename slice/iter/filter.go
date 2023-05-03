package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// FilterIter is the array based Iterator implementation that provides filtering of elements by a Predicate.
type FilterIter[T any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	filter   func(T) bool
}

var _ c.Iterator[any] = (*FilterIter[any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *FilterIter[T]) For(walker func(element T) error) error {
	return loop.For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f *FilterIter[T]) ForEach(walker func(element T)) {
	loop.ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f *FilterIter[T]) Next() (T, bool) {
	return nextFiltered(f.array, f.size, f.elemSize, f.filter, &f.i)
}

// Cap returns the iterator capacity
func (f *FilterIter[T]) Cap() int {
	return f.size
}
