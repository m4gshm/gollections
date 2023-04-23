package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// Fit is the Iterator wrapper that provides filtering of elements by a Predicate.
type Fit[T any] struct {
	next func() (T, bool)
	by   func(T) bool
}

var (
	_ c.Iterator[any] = (*Fit[any])(nil)
	_ c.Iterator[any] = Fit[any]{}
)

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f Fit[T]) Next() (element T, ok bool) {
	if next, by := f.next, f.by; next != nil && by != nil {
		element, ok = nextFiltered(next, by)
	}
	return element, ok
}

func nextFiltered[T any](next func() (T, bool), filter func(T) bool) (v T, ok bool) {
	return loop.First(next, filter)
}
