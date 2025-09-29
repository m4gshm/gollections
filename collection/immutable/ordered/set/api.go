// Package set provides [github.com/m4gshm/gollections/collection/ordered.Set] constructors and helpers
package set

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/seq"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) ordered.Set[T] {
	return ordered.NewSet(elements...)
}

// New instantiates Set and copies elements to it.
func New[T comparable](elements []T) ordered.Set[T] {
	return ordered.NewSet(elements...)
}

// FromSeq instantiates a set with elements retrieved by the seq.
func FromSeq[T comparable](seq seq.Seq[T]) ordered.Set[T] {
	return ordered.SetFromSeq(seq)
}

// Sort copy the specified set with sorted elements
func Sort[T comparable, f constraints.Ordered](s ordered.Set[T], by func(T) f) ordered.Set[T] {
	return collection.Sort(s, by)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[From, To comparable](set ordered.Set[From], converter func(From) To) seq.Seq[To] {
	return collection.Convert(set, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func Conv[From, To comparable](set ordered.Set[From], converter func(From) (To, error)) seq.SeqE[To] {
	return collection.Conv(set, converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](set ordered.Set[From], flattener func(From) []To) seq.Seq[To] {
	return collection.Flat(set, flattener)
}

// Flatt returns a errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](set ordered.Set[From], flattener func(From) ([]To, error)) seq.SeqE[To] {
	return collection.Flatt(set, flattener)
}
