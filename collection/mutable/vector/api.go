// Package vector provides mutable.Vector constructors and helpers
package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/seq"
)

// Of instantiates a vector with the specified elements
func Of[T any](elements ...T) *mutable.Vector[T] {
	return mutable.NewVector(elements...)
}

// Empty instantiates Vector with zero capacity.
func Empty[T any]() *mutable.Vector[T] {
	return NewCap[T](0)
}

// NewCap creates a vector with a predefined capacity.
func NewCap[T any](capacity int) *mutable.Vector[T] {
	return mutable.NewVectorCap[T](capacity)
}

// FromSeq creates a vector with elements retrieved by the seq.
func FromSeq[T any](seq seq.Seq[T]) *mutable.Vector[T] {
	return mutable.VectorFromSeq(seq)
}

// Sort sorts the specified vector in-place by a converter that thransforms an element to an Ordered (int, string and so on).
func Sort[T any, F constraints.Ordered](v *mutable.Vector[T], by func(T) F) *mutable.Vector[T] {
	return collection.Sort(v, by)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[From, To any](vector *mutable.Vector[From], converter func(From) To) seq.Seq[To] {
	return collection.Convert(vector, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func Conv[From, To comparable](vector *mutable.Vector[From], converter func(From) (To, error)) seq.SeqE[To] {
	return collection.Conv(vector, converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[From, To any](vector *mutable.Vector[From], flattener func(From) []To) seq.Seq[To] {
	return collection.Flat(vector, flattener)
}

// Flatt returns an errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](vector *mutable.Vector[From], flattener func(From) ([]To, error)) seq.SeqE[To] {
	return collection.Flatt(vector, flattener)
}
