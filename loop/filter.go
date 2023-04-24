package loop

import (
	"github.com/m4gshm/gollections/c"
)

// FitIter is the Iterator wrapper that provides filtering of elements by a Predicate.
type FitIter[T any] struct {
	next func() (T, bool)
	by   func(T) bool
}

var (
	_ c.Iterator[any] = (*FitIter[any])(nil)
	_ c.Iterator[any] = FitIter[any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f FitIter[T]) Next() (element T, ok bool) {
	if next, by := f.next, f.by; next != nil && by != nil {
		element, ok = nextFiltered(next, by)
	}
	return element, ok
}

func nextFiltered[T any](next func() (T, bool), filter func(T) bool) (v T, ok bool) {
	return First(next, filter)
}
