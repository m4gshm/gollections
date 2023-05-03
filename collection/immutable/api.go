// Package immutable provides immutable collection implementations
package immutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// NewMap instantiates Map populated by the 'elements' map key/values
func NewMap[K comparable, V any](elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// NewMapKV converts a slice of key/value pairs into a Map instance
func NewMapKV[K comparable, V any](elements []c.KV[K, V]) Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// NewSet instantiates Set and copies elements to it
func NewSet[T comparable](elements []T) Set[T] {
	internal := make(map[T]struct{}, len(elements))
	for _, e := range elements {
		internal[e] = struct{}{}
	}
	return WrapSet(internal)
}

// NewVector instantiates Vector populated by the 'elements' slice
func NewVector[T any](elements []T) Vector[T] {
	return WrapVector(slice.Clone(elements))
}
