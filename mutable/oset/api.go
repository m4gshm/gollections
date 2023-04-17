package oset

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
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

func Convert[From, To comparable](s *ordered.Set[From], by func(From) To) c.Pipe[To, []To] {
	h := *(s.Head())
	return it.NewPipe[To](it.Convert(h, h.Next, by))
}

func Flatt[From, To comparable](s *ordered.Set[From], by func(From) []To) c.Pipe[To, []To] {
	h := *(s.Head())
	f := it.Flatt(h, h.Next, by)
	return it.NewPipe[To](&f)
}
