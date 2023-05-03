// Package mutable provides implementations of mutable containers.
package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// NewMap instantiates Map populated by the 'elements' map key/values
func NewMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// NewMapCap instantiates Map with a predefined capacity.
func NewMapCap[K comparable, V any](capacity int) *Map[K, V] {
	return WrapMap(make(map[K]V, capacity))
}

// NewMapKV converts a slice of key/value pairs into a Map instance
func NewMapKV[K comparable, V any](elements []c.KV[K, V]) *Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// NewSet instantiates Set and copies elements to it
func NewSet[T comparable](elements []T) *Set[T] {
	internal := make(map[T]struct{}, len(elements))
	for _, e := range elements {
		internal[e] = struct{}{}
	}
	return WrapSet(internal)
}

// NewSetCap creates a set with a predefined capacity
func NewSetCap[T comparable](capacity int) *Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

// SetFromLoop creates a set with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func SetFromLoop[T comparable](next func() (T, bool)) *Set[T] {
	internal := map[T]struct{}{}
	for e, ok := next(); ok; e, ok = next() {
		internal[e] = struct{}{}
	}
	return WrapSet(internal)
}

// NewVector instantiates Vector populated by the 'elements' slice
func NewVector[T any](elements []T) *Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// NewVectorCap instantiates Vector with a predefined capacity
func NewVectorCap[T any](capacity int) *Vector[T] {
	return WrapVector(make([]T, 0, capacity))
}
