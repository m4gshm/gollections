// Package vector provides the Vector (ordered) implementation
package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/loop"
)

// Of instantiates Vector with predefined elements.
func Of[T any](elements ...T) immutable.Vector[T] {
	return immutable.NewVector(elements)
}

// New instantiates Vector and copies elements to it.
func New[T any](elements []T) immutable.Vector[T] {
	return immutable.NewVector(elements)
}

// From creates a Vector instance with elements obtained by passing an iterator.
func From[T any](elements c.Iterator[T]) immutable.Vector[T] {
	return immutable.WrapVector(loop.ToSlice(elements.Next))
}

// Sort instantiates Vector and puts sorted elements to it.
func Sort[t any, f constraints.Ordered](v immutable.Vector[t], by func(t) f) immutable.Vector[t] {
	return v.Sort(func(e1, e2 t) bool { return by(e1) < by(e2) })
}

func Convert[From, To any](v immutable.Vector[From], by func(From) To) c.Pipe[To, []To] {
	h := v.Head()
	return it.NewPipe[To](it.Convert(h, h.Next, by))
}

func Flatt[From, To any](v immutable.Vector[From], by func(From) []To) c.Pipe[To, []To] {
	b := v.Head()
	f := it.Flatt(b, b.Next, by)
	return it.NewPipe[To](&f)
}