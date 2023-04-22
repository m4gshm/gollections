package set

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/iterable/transform"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/mutable/ordered"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) *mutable.Set[T] {
	return mutable.NewSet(elements)
}

// From creates a Set instance with elements obtained by passing an iterator.
func From[T comparable](elements c.Iterator[T]) *mutable.Set[T] {
	return mutable.ToSet(elements)
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
func Sort[T comparable, F constraints.Ordered](s mutable.Set[T], by func(T) F) *ordered.Set[T] {
	return s.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func Convert[From, To comparable](collection *mutable.Set[From], converter func(From) To) c.Transform[To] {
	h := collection.Head()
	return transform.New(iter.Convert(h.Next, converter).Next)
}

// Flatt instantiates Iterator that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](s *mutable.Set[From], by func(From) []To) c.Transform[To] {
	h := s.Head()
	f := iter.Flatt(h.Next, by)
	return transform.New(f.Next)
}
