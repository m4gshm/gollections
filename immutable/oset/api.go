// Package oset provides the ordered set container implementation
package oset

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) ordered.Set[T] {
	return ordered.NewSet(elements)
}

// New instantiates Set and copies elements to it.
func New[T comparable](elements []T) ordered.Set[T] {
	return ordered.NewSet(elements)
}

// From creates a Set instance with elements obtained by passing an iterator.
func From[T comparable](elements c.Iterator[T]) ordered.Set[T] {
	return ordered.ToSet(elements)
}

// Sort instantiates Set and puts sorted elements to it.
func Sort[T comparable, f constraints.Ordered](s ordered.Set[T], by func(T) f) ordered.Set[T] {
	return s.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func Convert[From, To comparable](collection ordered.Set[From], converter func(From) To) c.Pipe[To] {
	h := collection.Head()
	return iter.NewPipe(iter.Convert(h.Next, converter).Next)
}

// Flatt returns a pipe that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](s ordered.Set[From], flattener func(From) []To) c.Pipe[To] {
	h := s.Head()
	f := iter.Flatt(h.Next, flattener)
	return iter.NewPipe(f.Next)
}
