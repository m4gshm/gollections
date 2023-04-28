// Package oset provides mutable ordered.Set constructors and helpers
package oset

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	breakLoop "github.com/m4gshm/gollections/loop/break/loop"
	"github.com/m4gshm/gollections/mutable/ordered"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.NewSet(elements)
}

// From instantiates a set with elements retrieved by the 'next' function
func From[T comparable](next func() (T, bool)) *ordered.Set[T] {
	return ordered.SetFromLoop(next)
}

// Empty instantiates Set with zero capacity.
func Empty[T comparable]() *ordered.Set[T] {
	return NewCap[T](0)
}

// NewCap instantiates Set with a predefined capacity.
func NewCap[T comparable](capacity int) *ordered.Set[T] {
	return ordered.NewSetCap[T](capacity)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To comparable](collection *ordered.Set[From], converter func(From) To) c.Stream[To] {
	h := collection.Head()
	return loop.Stream(loop.Convert(h.Next, converter).Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](collection *ordered.Set[From], converter func(From) (To, error)) c.StreamBreakable[To] {
	h := collection.Head()
	return breakLoop.Stream(breakLoop.Conv(breakLoop.From(h.Next), converter).Next)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](s *ordered.Set[From], flattener func(From) []To) c.Stream[To] {
	h := s.Head()
	f := loop.Flatt(h.Next, flattener)
	return loop.Stream(f.Next)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](s *ordered.Set[From], flattener func(From) ([]To, error)) c.StreamBreakable[To] {
	h := s.Head()
	f := breakLoop.Flat(breakLoop.From(h.Next), flattener)
	return breakLoop.Stream(f.Next)
}
