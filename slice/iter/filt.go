package iter

import (
	"unsafe"

	"github.com/m4gshm/gollections/break/c"
	"github.com/m4gshm/gollections/break/loop"
)

// FiltIter is the array based Iterator implementation that provides filtering of elements by a Predicate.
type FiltIter[T any] struct {
	array    unsafe.Pointer
	elemSize uintptr
	size, i  int
	filter   func(T) (bool, error)
}

var _ c.Iterator[any] = (*FiltIter[any])(nil)
var _ c.IterFor[any, *FiltIter[any]] = (*FiltIter[any])(nil)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f *FiltIter[T]) For(walker func(element T) error) error {
	return loop.For(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f *FiltIter[T]) Next() (T, bool, error) {
	return nextFilt(f.array, f.size, f.elemSize, f.filter, &f.i)
}

// Cap returns the iterator capacity
func (f *FiltIter[T]) Cap() int {
	return f.size
}

// Start is used with for loop construct like 'for i, val, ok, err := i.Start(); ok || err != nil ; val, ok, err = i.Next() { if err != nil { return err }}'
func (f *FiltIter[T]) Start() (*FiltIter[T], T, bool, error) {
	return startBreakIt[T](f)
}
