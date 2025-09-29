// Package immutable provides immutable collection implementations
package immutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) Set[T] {
	return SetFromSeq(seq.Of(elements...))
}

// NewSetOrdered instantiates ordered set and copies elements to it
func NewSetOrdered[T comparable](elements ...T) ordered.Set[T] {
	return ordered.NewSet(elements...)
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

// NewMap instantiates a map using key/value pairs
func NewMap[K comparable, V any](elements ...c.KV[K, V]) Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// NewMapOrdered instantiates an ordered map using key/value pairs
func NewMapOrdered[K comparable, V any](elements ...c.KV[K, V]) ordered.Map[K, V] {
	return ordered.NewMap(elements...)
}

// NewMapOf instantiates Map populated by the 'elements' map key/values
func NewMapOf[K comparable, V any](elements map[K]V) Map[K, V] {
	return MapFromSeq2(seq2.OfMap(elements))
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

// VectorFromSeq creates a vector with elements retrieved by the seq.
func VectorFromSeq[T any](s seq.Seq[T]) Vector[T] {
	return WrapVector(seq.Slice(s))
}
