package stream

import (
	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// / New instantiates Iter
func New[T any](next func() (T, bool)) Iter[T] {
	return Iter[T]{next: next}
}

// Iter is the Iterator based stream implementation.
type Iter[T any] struct {
	next func() (T, bool)
}

var (
	_ Stream[any]                                         = (*Iter[any])(nil)
	_ Stream[any]                                         = Iter[any]{}
	_ loop.Looper[any, Iter[any]]                         = (*Iter[any])(nil)
	_ loop.Looper[any, Iter[any]]                         = Iter[any]{}
	_ c.Filterable[any, Iter[any], breakStream.Iter[any]] = (*Iter[any])(nil)
	_ c.Filterable[any, Iter[any], breakStream.Iter[any]] = Iter[any]{}
	_ c.Iterator[any]                                     = (*Iter[any])(nil)
	_ c.Iterator[any]                                     = Iter[any]{}
)

// Next implements c.Iterator
func (t Iter[T]) Next() (element T, ok bool) {
	if next := t.next; next != nil {
		element, ok = next()
	}
	return element, ok
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t Iter[T]) Filter(predicate func(T) bool) Iter[T] {
	f := loop.Filter(t.next, predicate)
	return New(f.Next)
}

// Filt returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (t Iter[T]) Filt(predicate func(T) (bool, error)) breakStream.Iter[T] {
	f := breakLoop.Filt(breakLoop.From(t.next), predicate)
	return breakStream.New(f.Next)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (t Iter[T]) Convert(converter func(T) T) Iter[T] {
	conv := loop.Convert(t.next, converter)
	return New(conv.Next)
}

// Conv returns a stream that applies the 'converter' function to the collection elements
func (t Iter[T]) Conv(converter func(T) (T, error)) breakStream.Iter[T] {
	conv := breakLoop.Conv(breakLoop.From(t.next), converter)
	return breakStream.New(conv.Next)
}

// ForEach applies the 'walker' function for every element
func (t Iter[T]) ForEach(walker func(T)) {
	loop.ForEach(t.next, walker)
}

// For applies the 'walker' function for the elements. Return the c.ErrBreak to stop.
func (t Iter[T]) For(walker func(T) error) error {
	return loop.For(t.next, walker)
}

// Reduce reduces the elements into an one using the 'merge' function
func (t Iter[T]) Reduce(by func(T, T) T) T {
	return loop.Reduce(t.next, by)
}

// Iter returns as an iterator
func (t Iter[T]) Iter() c.Iterator[T] {
	return t
}

func (t Iter[T]) Loop() Iter[T] {
	return t
}

// Slice collects the elements to a slice
func (t Iter[T]) Slice() []T {
	return loop.ToSlice(t.next)
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func (t Iter[T]) HasAny(predicate func(T) bool) bool {
	return loop.HasAny(t.next, predicate)
}
