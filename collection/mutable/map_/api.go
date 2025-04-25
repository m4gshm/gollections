// Package map_ provides unordered mutable.Map constructors
package map_

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/mutable"
	"github.com/m4gshm/gollections/seq"
)

// Of instantiates a ap from the specified key/value pairs
func Of[K comparable, V any](elements ...c.KV[K, V]) *mutable.Map[K, V] {
	return mutable.NewMap(elements...)
}

// Empty instantiates a Map with zero capacity.
func Empty[K comparable, V any]() *mutable.Map[K, V] {
	return New[K, V](0)
}

// New instantiates a Map with a predefined capacity.
func New[K comparable, V any](capacity int) *mutable.Map[K, V] {
	return mutable.NewMapCap[K, V](capacity)
}

// From instantiates a map with elements obtained by passing the 'loop' function.
// Deprecated: replaced by [FromSeq2].
func From[K comparable, V any](next func() (K, V, bool)) *mutable.Map[K, V] {
	return mutable.MapFromLoop(next)
}

// FromSeq2 creates a map with elements retrieved by the seq.
func FromSeq2[K comparable, V any](seq seq.Seq2[K, V]) *mutable.Map[K, V] {
	return mutable.MapFromSeq2(seq)
}
