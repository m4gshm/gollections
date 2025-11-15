// Package set provides unordered [github.com/m4gshm/gollections/collection/immutable.Set] constructors and helpers
package set

import (
	"cmp"

	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/seq"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) immutable.Set[T] {
	return immutable.NewSet(elements...)
}

// New instantiates Set and copies elements to it.
func New[T comparable](elements []T) immutable.Set[T] {
	return immutable.NewSet(elements...)
}

// FromSeq creates a set with elements retrieved by the seq.
func FromSeq[T comparable](seq seq.Seq[T]) immutable.Set[T] {
	return immutable.SetFromSeq(seq)
}

// Sort instantiates Set and puts sorted elements to it.
func Sort[T comparable, f cmp.Ordered](s immutable.Set[T], by func(T) f) ordered.Set[T] {
	return collection.Sort(s, by)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func Convert[From, To comparable](set immutable.Set[From], converter func(From) To) seq.Seq[To] {
	return collection.Convert(set, converter)
}

// Conv returns an errorable seq that applies the 'converter' function to the collection elements
func Conv[From, To comparable](set immutable.Set[From], converter func(From) (To, error)) seq.SeqE[To] {
	return collection.Conv(set, converter)
}

// Flat returns a seq that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](set immutable.Set[From], flattener func(From) []To) seq.Seq[To] {
	return collection.Flat(set, flattener)
}

// Flatt returns an errorable seq that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](set immutable.Set[From], flattener func(From) ([]To, error)) seq.SeqE[To] {
	return collection.Flatt(set, flattener)
}
