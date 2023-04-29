// Package set provides unordered immutable.Set constructors and helpers
package set

import (
	"golang.org/x/exp/constraints"

	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/stream"
)

// Of instantiates Set with predefined elements.
func Of[T comparable](elements ...T) immutable.Set[T] {
	return immutable.NewSet(elements)
}

// New instantiates Set and copies elements to it.
func New[T comparable](elements []T) immutable.Set[T] {
	return immutable.NewSet(elements)
}

// From instantiates a map with key/values retrieved by the 'next' function.
// The next returns a key/value pairs with true or zero values with false if there are no more elements.
func From[T comparable](next func() (T, bool)) immutable.Set[T] {
	return immutable.SetFromLoop(next)
}

// Sort instantiates Set and puts sorted elements to it.
func Sort[T comparable, f constraints.Ordered](s immutable.Set[T], by func(T) f) ordered.Set[T] {
	return s.Sort(func(e1, e2 T) bool { return by(e1) < by(e2) })
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func Convert[From, To comparable](collection immutable.Set[From], converter func(From) To) stream.Iter[To] {
	h := collection.Head()
	return stream.New(loop.Convert(h.Next, converter).Next)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func Conv[From, To comparable](collection immutable.Set[From], converter func(From) (To, error)) breakStream.Iter[To] {
	h := collection.Head()
	return breakStream.New(breakLoop.Conv(breakLoop.From(h.Next), converter).Next)
}

// Flatt returns a stream that converts the collection elements into slices and then flattens them to one level
func Flatt[From, To comparable](s immutable.Set[From], flattener func(From) []To) stream.Iter[To] {
	h := s.Head()
	f := loop.Flatt(h.Next, flattener)
	return stream.New(f.Next)
}

// Flat returns a breakable stream that converts the collection elements into slices and then flattens them to one level
func Flat[From, To comparable](s immutable.Set[From], flattener func(From) ([]To, error)) breakStream.Iter[To] {
	h := s.Head()
	f := breakLoop.Flat(breakLoop.From(h.Next), flattener)
	return breakStream.New(f.Next)
}
