// Package stream provides a stream implementation and helper functions
package stream

import (
	"github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/convert/as"
)

// New instantiates a stream instance
func New[T any](next func() (T, bool, error)) Iter[T] {
	return Iter[T]{next: next}
}

// Iter is the Iterator based stream implementation.
type Iter[T any] struct {
	next func() (T, bool, error)
}

var (
	_ Stream[any, Iter[any]]                  = (*Iter[any])(nil)
	_ Stream[any, Iter[any]]                  = Iter[any]{}
	_ c.Filterable[any, Iter[any], Iter[any]] = Iter[any]{}
)

// Next implements c.Iterator
func (t Iter[T]) Next() (element T, ok bool, err error) {
	if next := t.next; next != nil {
		element, ok, err = next()
	}
	return element, ok, err
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (t Iter[T]) Filt(predicate func(T) (bool, error)) Iter[T] {
	f := loop.Filt(t.next, predicate)
	return New(f.Next)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t Iter[T]) Filter(predicate func(T) bool) Iter[T] {
	f := loop.Filt(t.next, as.ErrTail(predicate))
	return New(f.Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (t Iter[T]) Convert(converter func(T) T) Iter[T] {
	conv := loop.Conv(t.next, as.ErrTail(converter))
	return New(conv.Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (t Iter[T]) Conv(converter func(T) (T, error)) Iter[T] {
	conv := loop.Conv(t.next, converter)
	return New(conv.Next)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (t Iter[T]) For(walker func(T) error) error {
	return loop.For(t.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (t Iter[T]) Reduce(merger func(T, T) (T, error)) (T, error) {
	return loop.Reducee(t.next, merger)
}

// First returns the first element that satisfies the condition of the 'predicate' function
func (t Iter[T]) First(predicate func(T) (bool, error)) (T, bool, error) {
	return loop.Firstt(t.next, predicate)
}

// Iter creates an iterator and returns as interface
func (t Iter[T]) Iter() Iter[T] {
	return t
}

// Slice collects the elements to a slice
func (t Iter[T]) Slice() ([]T, error) {
	return loop.Slice(t.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (t Iter[T]) HasAny(predicate func(T) (bool, error)) (bool, error) {
	return loop.HasAnyy(t.next, predicate)
}
