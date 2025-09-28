// Package set provides unordered [github.com/m4gshm/gollections/collection/mutable.Set] constructors and helpers
package set

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) *mutable.Set[T] {
	return mutable.NewSet(elements...)
}

// From instantiates a set with elements retrieved by the 'next' function.
//
// Deprecated: replaced by [FromSeq].
func From[T comparable](next func() (T, bool)) *mutable.Set[T] {
	return mutable.SetFromLoop(next)
}

// FromSeq creates a set with elements retrieved by the seq.
func FromSeq[T comparable](seq collection.Seq[T]) *mutable.Set[T] {
	return mutable.SetFromSeq(seq)
}

// Empty instantiates Set with zero capacity.
func Empty[T comparable]() *mutable.Set[T] {
	return NewCap[T](0)
}

// NewCap instantiates Set with a predefined capacity.
func NewCap[T comparable](capacity int) *mutable.Set[T] {
	return mutable.NewSetCap[T](capacity)
}

// Sort sorts a Set in-place by a converter that thransforms an element to an Ordered (int, string and so on).
func Sort[T comparable, F constraints.Ordered](s *mutable.Set[T], by func(T) F) *ordered.Set[T] {
	return collection.Sort(s, by)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[From, To comparable](set *mutable.Set[From], converter func(From) To) collection.Seq[To] {
	return collection.Convert(set, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func Conv[From, To comparable](set *mutable.Set[From], converter func(From) (To, error)) collection.SeqE[To] {
	return collection.Conv(set, converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](set *mutable.Set[From], flattener func(From) []To) collection.Seq[To] {
	return collection.Flat(set, flattener)
}

// Flatt returns a errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](set *mutable.Set[From], flattener func(From) ([]To, error)) collection.SeqE[To] {
	return collection.Flatt(set, flattener)
}
