// Package oset provides the ordered set container implementation
package oset

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/it/impl/it"
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

func Convert[From, To comparable](s ordered.Set[From], by func(From) To) c.Pipe[To, []To] {
	h := s.Head()
	return it.NewPipe[To](it.Convert(h, h.Next, by))
}

func Flatt[From, To comparable](s ordered.Set[From], by func(From) []To) c.Pipe[To, []To] {
	b := s.Head()
	f := it.Flatt(b, b.Next, by)
	return it.NewPipe[To](&f)
}
