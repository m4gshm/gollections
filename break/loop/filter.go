package loop

import (
	"github.com/m4gshm/gollections/break/c"
)

// FiltIter is the Iterator wrapper that provides filtering of elements by a Predicate.
type FiltIter[T any] struct {
	next   func() (T, bool, error)
	filter func(T) (bool, error)
}

var (
	_ c.Iterator[any] = (*FiltIter[any])(nil)
	_ c.Iterator[any] = FiltIter[any]{}
)

var _ c.IterFor[any, FiltIter[any]] = (*FiltIter[any])(nil)

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
	return element, ok, err
}

// Start is used with for loop construct like 'for i, val, ok, err := i.Start(); ok || err != nil ; val, ok, err = i.Next() { if err != nil { return err }}'
func (f FiltIter[T]) Start() (FiltIter[T], T, bool, error) {
	return startIt[T](f)
}

func nextFiltered[T any](next func() (T, bool, error), filter func(T) (bool, error)) (v T, ok bool, err error) {
	return Firstt(next, filter)
}
