// Package immutable provides immutable collection implementations
package immutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) Set[T] {
	uniques := make(map[T]struct{}, len(elements))
	for _, e := range elements {
		uniques[e] = struct{}{}
	}
	return WrapSet(uniques)
}

// SetFromLoop creates a set with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func SetFromLoop[T comparable](next func() (T, bool)) Set[T] {
	uniques := map[T]struct{}{}
	for e, ok := next(); ok; e, ok = next() {
		uniques[e] = struct{}{}
	}
	return WrapSet(uniques)
}

// NewMap instantiates an map using key/value pairs
func NewMap[K comparable, V any](elements ...c.KV[K, V]) Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// NewMapOf instantiates Map populated by the 'elements' map key/values
func NewMapOf[K comparable, V any](elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// MapFromLoop creates a map with elements retrieved converter the 'next' function
func MapFromLoop[K comparable, V any](next func() (K, V, bool)) Map[K, V] {
	uniques := map[K]V{}
	for key, val, ok := next(); ok; key, val, ok = next() {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// NewVector instantiates Vector populated by the 'elements' slice
func NewVector[T any](elements ...T) Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// VectorFromLoop creates a vector with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func VectorFromLoop[T any](next func() (T, bool)) Vector[T] {
	return WrapVector(loop.Slice(next))
}
