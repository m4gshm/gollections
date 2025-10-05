// Package vector provides ordered immutable.Vector constructors and helpers
package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/seq"
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

// Sort copy the specified vector with sorted elements
func Sort[T any, F constraints.Ordered](v immutable.Vector[T], by func(T) F) immutable.Vector[T] {
	return collection.Sort(v, by)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[From, To any](vector immutable.Vector[From], converter func(From) To) seq.Seq[To] {
	return collection.Convert(vector, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func Conv[From, To comparable](vector immutable.Vector[From], converter func(From) (To, error)) seq.SeqE[To] {
	return collection.Conv(vector, converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[From any, To any](vector immutable.Vector[From], flattener func(From) []To) seq.Seq[To] {
	return collection.Flat(vector, flattener)
}

// Flatt returns an errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](vector immutable.Vector[From], flattener func(From) ([]To, error)) seq.SeqE[To] {
	return collection.Flatt(vector, flattener)
}
