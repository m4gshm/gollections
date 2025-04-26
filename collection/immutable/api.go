// Package immutable provides immutable collection implementations
package immutable

import (
	"github.com/m4gshm/gollections/c"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) Set[T] {
	return SetFromSeq(seq.Of(elements...))
}

// SetFromLoop creates a set with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
// Deprecated: replaced by [SetFromSeq].
func SetFromLoop[T comparable](next func() (T, bool)) Set[T] {
	return SetFromSeq((loop.Loop[T])(next).All)
}

// SetFromSeq creates a set with elements retrieved by the seq.
func SetFromSeq[T comparable](seq seq.Seq[T]) Set[T] {
	if seq == nil {
		return Set[T]{}
	}
	uniques := map[T]struct{}{}
	for e := range seq {
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
	return MapFromSeq2(seq2.OfMap(elements))
}

// MapFromLoop creates a map with elements retrieved converter the 'next' function.
// Deprecated: replaced by [MapFromSeq2].
func MapFromLoop[K comparable, V any](next func() (K, V, bool)) Map[K, V] {
	return MapFromSeq2(kvloop.Loop[K, V](next).All)
}

// MapFromSeq2 creates a map with elements retrieved by the seq.
func MapFromSeq2[K comparable, V any](seq seq.Seq2[K, V]) Map[K, V] {
	uniques := map[K]V{}
	for key, val := range seq {
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
// Deprecated: replaced by [VectorFromLoop].
func VectorFromLoop[T any](next func() (T, bool)) Vector[T] {
	return WrapVector(loop.Slice(next))
}

// VectorFromSeq creates a vector with elements retrieved by the seq.
func VectorFromSeq[T any](s seq.Seq[T]) Vector[T] {
	return WrapVector(seq.Slice(s))
}
