// Package map_ provides immutable.Map constructors
package map_ // nilint

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/seq"
)

// Of instantiates a ap from the specified key/value pairs
func Of[K comparable, V any](elements ...c.KV[K, V]) immutable.Map[K, V] {
	return immutable.NewMap(elements...)
}

// New instantiates Map and copies elements to it
func New[K comparable, V any](elements map[K]V) immutable.Map[K, V] {
	return immutable.NewMapOf(elements)
}

// FromSeq2 creates a map with elements retrieved by the seq.
func FromSeq2[K comparable, V any](seq seq.Seq2[K, V]) immutable.Map[K, V] {
	return immutable.MapFromSeq2(seq)
}
