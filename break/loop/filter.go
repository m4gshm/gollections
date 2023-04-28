package loop

import (
	"github.com/m4gshm/gollections/c"
)

// FiltIter is the Iterator wrapper that provides filtering of elements by a Predicate.
type FiltIter[T any] struct {
	next   func() (T, bool, error)
	filter func(T) (bool, error)
}

var (
	_ c.IteratorBreakable[any] = (*FiltIter[any])(nil)
	_ c.IteratorBreakable[any] = FiltIter[any]{}
)

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f FiltIter[T]) For(walker func(element T) error) error {
	return For(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f FiltIter[T]) Next() (element T, ok bool, err error) {
	if next, by := f.next, f.filter; next != nil && by != nil {
		element, ok, err = nextFiltered(next, by)
	}
	return element, ok, nil
}

func nextFiltered[T any](next func() (T, bool, error), filter func(T) (bool, error)) (v T, ok bool, err error) {
	return First(next, filter)
}
