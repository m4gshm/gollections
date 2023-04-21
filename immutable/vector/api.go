// Package vector provides the Vector (ordered) implementation
package vector

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/iter/impl/iter"
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

// Convert returns a pipe that applies the 'converter' function to the collection elements
func Convert[From, To any](collection immutable.Vector[From], converter func(From) To) c.Pipe[To] {
	h := collection.Head()
	return iter.NewPipe[To](iter.Convert(h, h.Next, converter))
}

// Flatt returns a pipe that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To any](collection immutable.Vector[From], by func(From) []To) c.Pipe[To] {
	h := collection.Head()
	f := iter.Flatt(h, h.Next, by)
	return iter.NewPipe[To](&f)
}
