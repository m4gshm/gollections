package oset

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/mutable/ordered"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.NewSet(elements)
}

// From creates a Set instance with elements obtained by passing an iterator.
func From[T comparable](elements c.Iterator[T]) *ordered.Set[T] {
	return ordered.ToSet(elements)
}

// Empty instantiates Set with zero capacity.
func Empty[T comparable]() *ordered.Set[T] {
	return NewCap[T](0)
}

// NewCap instantiates Set with a predefined capacity.
func NewCap[T comparable](capacity int) *ordered.Set[T] {
	return ordered.NewSetCap[T](capacity)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func Convert[From, To comparable](collection *ordered.Set[From], converter func(From) To) c.Pipe[To] {
	h := collection.Head()
	return iter.NewPipe[To](iter.Convert(h, h.Next, converter))
}

// Flatt instantiates Iterator that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](s *ordered.Set[From], by func(From) []To) c.Pipe[To] {
	h := s.Head()
	f := iter.Flatt(h, h.Next, by)
	return iter.NewPipe[To](&f)
}
