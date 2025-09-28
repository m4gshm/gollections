// Package vector provides ordered immutable.Vector constructors and helpers
package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
)

// Of instantiates a vector with the specified elements
func Of[T any](elements ...T) immutable.Vector[T] {
	return immutable.NewVector(elements...)
}

// New instantiates a vector with the specified elements
func New[T any](elements []T) immutable.Vector[T] {
	return immutable.NewVector(elements...)
}

// Wrap instantiates Vector using a slise as internal storage.
func Wrap[T any](elements []T) immutable.Vector[T] {
	return immutable.WrapVector(elements)
}

// From instantiates a vector with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func From[T any](next func() (T, bool)) immutable.Vector[T] {
	return immutable.VectorFromLoop(next)
}

// Sort copy the specified vector with sorted elements
func Sort[T any, F constraints.Ordered](v immutable.Vector[T], by func(T) F) immutable.Vector[T] {
	return collection.Sort[immutable.Vector[T]](v, by)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[From, To any](vector immutable.Vector[From], converter func(From) To) collection.Seq[To] {
	return collection.Convert(vector, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func Conv[From, To comparable](vector immutable.Vector[From], converter func(From) (To, error)) collection.SeqE[To] {
	return collection.Conv(vector, converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[From any, To any](vector immutable.Vector[From], flattener func(From) []To) collection.Seq[To] {
	return collection.Flat(vector, flattener)
}

// Flatt returns a errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](vector immutable.Vector[From], flattener func(From) ([]To, error)) collection.SeqE[To] {
	return collection.Flatt(vector, flattener)
}
