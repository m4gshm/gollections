// Package oset provides ordered.Set constructors and helpers
package oset

import (
	"golang.org/x/exp/constraints"

	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/immutable/ordered"
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
func Convert[From, To comparable](set ordered.Set[From], converter func(From) To) stream.Iter[To] {
	return collection.Convert(set, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](set ordered.Set[From], converter func(From) (To, error)) breakStream.Iter[To] {
	return collection.Conv(set, converter)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](set ordered.Set[From], flattener func(From) []To) stream.Iter[To] {
	return collection.Flatt(set, flattener)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](set ordered.Set[From], flattener func(From) ([]To, error)) breakStream.Iter[To] {
	return collection.Flat(set, flattener)
}
