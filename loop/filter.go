package loop

import (
	"github.com/m4gshm/gollections/c"
)

// FiltIter is the Iterator wrapper that provides filtering of elements by a Predicate.
type FiltIter[T any] struct {
	next func() (T, bool)
	by   func(T) bool
}

var (
	_ c.Iterator[any] = (*FiltIter[any])(nil)
	_ c.Iterator[any] = FiltIter[any]{}
)

var _ c.IterFor[any, FiltIter[any]] = FiltIter[any]{}

// For takes elements retrieved by the iterator. Can be interrupt by returning ErrBreak
func (f FiltIter[T]) For(walker func(element T) error) error {
	return For(f.Next, walker)
}

// ForEach FlatIter all elements retrieved by the iterator
func (f FiltIter[T]) ForEach(walker func(element T)) {
	ForEach(f.Next, walker)
}

// Next returns the next element.
// The ok result indicates whether the element was returned by the iterator.
// If ok == false, then the iteration must be completed.
func (f FiltIter[T]) Next() (element T, ok bool) {
	if next, by := f.next, f.by; next != nil && by != nil {
		element, ok = nextFiltered(next, by)
	}
	return element, ok
}

// Start is used with for loop construct like 'for i, val, ok := i.Start(); ok; val, ok = i.Next() { }'
func (f FiltIter[T]) Start() (FiltIter[T], T, bool) {
	return startIt[T](f)
}

func nextFiltered[T any](next func() (T, bool), filter func(T) bool) (v T, ok bool) {
	return First(next, filter)
}
