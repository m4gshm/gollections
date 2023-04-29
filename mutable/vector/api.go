// Package vector provides mutable.Vector constructors and helpers
package vector

import (
	"golang.org/x/exp/constraints"

	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/stream"
)

// Of instantiates a vector with the specified elements
func Of[T any](elements ...T) *mutable.Vector[T] {
	return mutable.NewVector(elements)
}

// Empty instantiates Vector with zero capacity.
func Empty[T any]() *mutable.Vector[T] {
	return NewCap[T](0)
}

// NewCap creates a vector with a predefined capacity.
func NewCap[T any](capacity int) *mutable.Vector[T] {
	return mutable.NewVectorCap[T](capacity)
}

// From instantiates a vector with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func From[T any](next func() (T, bool)) *mutable.Vector[T] {
	return mutable.WrapVector(loop.Slice(next))
}

// Sort sorts the specified vector in-place by a converter that thransforms an element to an Ordered (int, string and so on).
func Sort[T any, F constraints.Ordered](v *mutable.Vector[T], by func(T) F) *mutable.Vector[T] {
	return v.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To any](collection *mutable.Vector[From], converter func(From) To) stream.Iter[To] {
	return iterable.Convert(collection, converter)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](collection *mutable.Vector[From], converter func(From) (To, error)) breakStream.Iter[To] {
	return iterable.Conv(collection, converter)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To any](collection *mutable.Vector[From], flattener func(From) []To) stream.Iter[To] {
	return iterable.Flatt(collection, flattener)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](collection *mutable.Vector[From], flattener func(From) ([]To, error)) breakStream.Iter[To] {
	return iterable.Flat(collection, flattener)
}
