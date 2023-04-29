// Package oset provides mutable ordered.Set constructors and helpers
package oset

import (
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/mutable/ordered"
	"github.com/m4gshm/gollections/stream"
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
func Convert[From, To comparable](collection *ordered.Set[From], converter func(From) To) stream.Iter[To] {
	return iterable.Convert(collection, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](collection *ordered.Set[From], converter func(From) (To, error)) breakStream.Iter[To] {
	return iterable.Conv(collection, converter)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](collection *ordered.Set[From], flattener func(From) []To) stream.Iter[To] {
	return iterable.Flatt(collection, flattener)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](collection *ordered.Set[From], flattener func(From) ([]To, error)) breakStream.Iter[To] {
	return iterable.Flat(collection, flattener)
}
