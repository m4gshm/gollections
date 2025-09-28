// Package map_ provides mutable ordered.Map constructors
package map_

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
)

// Of instantiates a ap from the specified key/value pairs
func Of[K comparable, V any](pairs ...c.KV[K, V]) *ordered.Map[K, V] {
	return ordered.NewMap(pairs...)
}

// Empty instantiates a map with zero capacity
func Empty[K comparable, V any]() *ordered.Map[K, V] {
	return New[K, V](0)
}

// New instantiates a map with a predefined capacity
func New[K comparable, V any](capacity int) *ordered.Map[K, V] {
	return ordered.WrapMap(make([]K, 0, capacity), make(map[K]V, capacity))
}

// From instantiates a map with elements obtained by passing the 'loop' function.
//
// Deprecated: replaced by [MapFromSeq2].
func From[K comparable, V any](next func() (K, V, bool)) *ordered.Map[K, V] {
	return ordered.MapFromLoop(next)
}

// FromSeq2 creates a map with elements retrieved by the seq.
func FromSeq2[K comparable, V any](seq collection.Seq2[K, V]) *ordered.Map[K, V] {
	return ordered.MapFromSeq2(seq)
}
