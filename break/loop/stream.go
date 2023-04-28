package loop

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
)

// / Stream is default Stream constructor/
func Stream[T any](next func() (T, bool, error)) StreamIter[T] {
	return StreamIter[T]{next: next}
}

// StreamIter is the Iterator based stream implementation.
type StreamIter[T any] struct {
	next func() (T, bool, error)
}

var (
	_ c.StreamBreakable[any, StreamIter[any]]             = (*StreamIter[any])(nil)
	_ c.StreamBreakable[any, StreamIter[any]]             = StreamIter[any]{}
	_ c.Filterable[any, StreamIter[any], StreamIter[any]] = StreamIter[any]{}
)

// Next implements c.Iterator
func (t StreamIter[T]) Next() (element T, ok bool, err error) {
	if next := t.next; next != nil {
		element, ok, err = next()
	}
	return element, ok, err
}

// Filt returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t StreamIter[T]) Filt(predicate func(T) (bool, error)) StreamIter[T] {
	f := Filt(t.next, predicate)
	return Stream(f.Next)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t StreamIter[T]) Filter(predicate func(T) bool) StreamIter[T] {
	f := Filt(t.next, as.ErrTail(predicate))
	return Stream(f.Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (t StreamIter[T]) Convert(converter func(T) T) StreamIter[T] {
	conv := Conv(t.next, as.ErrTail(converter))
	return Stream(conv.Next)
}

// Conv returns a stream that applies the 'converter' function to the collection elements
func (t StreamIter[T]) Conv(converter func(T) (T, error)) StreamIter[T] {
	conv := Conv(t.next, converter)
	return Stream(conv.Next)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (t StreamIter[T]) For(walker func(T) error) error {
	return For(t.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (t StreamIter[T]) Reduce(merger func(T, T) (T, error)) (T, error) {
	return Reduce(t.next, merger)
}

// Begin creates iterator
func (t StreamIter[T]) Begin() StreamIter[T] {
	return t
}

// Slice collects the elements to a slice
func (t StreamIter[T]) Slice() ([]T, error) {
	return ToSlice(t.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (t StreamIter[T]) HasAny(predicate func(T) (bool, error)) (bool, error) {
	return HasAny(t.next, predicate)
}
