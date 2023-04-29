// Package oset provides ordered.Set constructors and helpers
package oset

import (
	"golang.org/x/exp/constraints"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/stream"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) ordered.Set[T] {
	return ordered.NewSet(elements)
}

// New instantiates Set and copies elements to it.
func New[T comparable](elements []T) ordered.Set[T] {
	return ordered.NewSet(elements)
}

// From instantiates a set with elements retrieved by the 'next' function
func From[T comparable](next func() (T, bool)) ordered.Set[T] {
	return ordered.SetFromLoop(next)
}

// Sort copy the specified set with sorted elements
func Sort[T comparable, f constraints.Ordered](s ordered.Set[T], by func(T) f) ordered.Set[T] {
	return s.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To comparable](collection ordered.Set[From], converter func(From) To) stream.Iter[To] {
	h := collection.Head()
	return stream.New(loop.Convert(h.Next, converter).Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](collection ordered.Set[From], converter func(From) (To, error)) breakStream.Iter[To] {
	h := collection.Head()
	return breakStream.New(breakLoop.Conv(breakLoop.From(h.Next), converter).Next)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](s ordered.Set[From], flattener func(From) []To) stream.Iter[To] {
	h := s.Head()
	f := loop.Flatt(h.Next, flattener)
	return stream.New(f.Next)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](s ordered.Set[From], flattener func(From) ([]To, error)) breakStream.Iter[To] {
	h := s.Head()
	f := breakLoop.Flat(breakLoop.From(h.Next), flattener)
	return breakStream.New(f.Next)
}
