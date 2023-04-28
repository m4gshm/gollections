package loop

import (
	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
)

// / Stream is default Stream constructor/
func Stream[T any](next func() (T, bool)) StreamIter[T] {
	return StreamIter[T]{next: next}
}

// StreamIter is the Iterator based stream implementation.
type StreamIter[T any] struct {
	next func() (T, bool)
}

var (
	_ c.Stream[any, StreamIter[any]]                                = (*StreamIter[any])(nil)
	_ c.Stream[any, StreamIter[any]]                                = StreamIter[any]{}
	_ c.Filterable[any, StreamIter[any], breakLoop.StreamIter[any]] = (*StreamIter[any])(nil)
	_ c.Filterable[any, StreamIter[any], breakLoop.StreamIter[any]] = StreamIter[any]{}
	_ c.Iterator[any]                                               = (*StreamIter[any])(nil)
	_ c.Iterator[any]                                               = StreamIter[any]{}
)

// Next implements c.Iterator
func (t StreamIter[T]) Next() (element T, ok bool) {
	if next := t.next; next != nil {
		element, ok = next()
	}
	return element, ok
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t StreamIter[T]) Filter(predicate func(T) bool) StreamIter[T] {
	f := Filter(t.next, predicate)
	return Stream(f.Next)
}

// Filt returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t StreamIter[T]) Filt(predicate func(T) (bool, error)) breakLoop.StreamIter[T] {
	f := breakLoop.Filt(breakLoop.From(t.next), predicate)
	return breakLoop.Stream(f.Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (t StreamIter[T]) Convert(converter func(T) T) StreamIter[T] {
	conv := Convert(t.next, converter)
	return Stream(conv.Next)
}

// Conv returns a stream that applies the 'converter' function to the collection elements
func (t StreamIter[T]) Conv(converter func(T) (T, error)) breakLoop.StreamIter[T] {
	conv := breakLoop.Conv(breakLoop.From(t.next), converter)
	return breakLoop.Stream(conv.Next)
}

// ForEach applies the 'walker' function for every element
func (t StreamIter[T]) ForEach(walker func(T)) {
	ForEach(t.next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (t StreamIter[T]) For(walker func(T) error) error {
	return For(t.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (t StreamIter[T]) Reduce(by func(T, T) T) T {
	return Reduce(t.next, by)
}

// Begin creates iterator
func (t StreamIter[T]) Begin() StreamIter[T] {
	return t
}

// Slice collects the elements to a slice
func (t StreamIter[T]) Slice() []T {
	return ToSlice(t.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (t StreamIter[T]) HasAny(predicate func(T) bool) bool {
	return HasAny(t.next, predicate)
}
