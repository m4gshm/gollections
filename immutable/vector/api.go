// Package vector provides ordered immutable.Vector constructors and helpers
package vector

import (
	"golang.org/x/exp/constraints"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/stream"
)

// Of instantiates a vector with the specified elements
func Of[T any](elements ...T) immutable.Vector[T] {
	return immutable.NewVector(elements)
}

// New instantiates a vector with the specified elements
func New[T any](elements []T) immutable.Vector[T] {
	return immutable.NewVector(elements)
}

// From instantiates a vector with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func From[T any](next func() (T, bool)) immutable.Vector[T] {
	return immutable.WrapVector(loop.ToSlice(next))
}

// Sort copy the specified vector with sorted elements
func Sort[t any, f constraints.Ordered](v immutable.Vector[t], by func(t) f) immutable.Vector[t] {
	return v.Sort(func(e1, e2 t) bool { return by(e1) < by(e2) })
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To any](collection immutable.Vector[From], converter func(From) To) stream.Iter[To] {
	h := collection.Head()
	return stream.New(loop.Convert(h.Next, converter).Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](collection immutable.Vector[From], converter func(From) (To, error)) breakStream.Iter[To] {
	h := collection.Head()
	return breakStream.New(breakLoop.Conv(breakLoop.From(h.Next), converter).Next)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From any, To any](collection immutable.Vector[From], flattener func(From) []To) stream.Iter[To] {
	h := collection.Head()
	f := loop.Flatt(h.Next, flattener)
	return stream.New(f.Next)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](s immutable.Vector[From], flattener func(From) ([]To, error)) breakStream.Iter[To] {
	h := s.Head()
	f := breakLoop.Flat(breakLoop.From(h.Next), flattener)
	return breakStream.New(f.Next)
}
