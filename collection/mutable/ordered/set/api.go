// Package set provides mutable [github.com/m4gshm/gollections/collection/ordered.Set] constructors and helpers
package set

import (
	"golang.org/x/exp/constraints"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/loop"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.NewSet(elements...)
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

// Sort copy the specified set with sorted elements
func Sort[T comparable, O constraints.Ordered](s *ordered.Set[T], by func(T) O) *ordered.Set[T] {
	return collection.Sort[*ordered.Set[T]](s, by)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func Convert[From, To comparable](set *ordered.Set[From], converter func(From) To) loop.Loop[To] {
	return collection.Convert(set, converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func Conv[From, To comparable](set *ordered.Set[From], converter func(From) (To, error)) breakLoop.Loop[To] {
	return collection.Conv(set, converter)
}

// Flat returns a loop that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](set *ordered.Set[From], flattener func(From) []To) loop.Loop[To] {
	return collection.Flat(set, flattener)
}

// Flatt returns a breakable loop that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](set *ordered.Set[From], flattener func(From) ([]To, error)) breakLoop.Loop[To] {
	return collection.Flatt(set, flattener)
}
