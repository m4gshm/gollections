// Package mutable provides implementations of mutable containers.
package mutable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) *Set[T] {
	return SetFromSeq(seq.Of(elements...))
}

// NewSetOrdered instantiates ordered set and copies elements to it
func NewSetOrdered[T comparable](elements ...T) *ordered.Set[T] {
	return ordered.NewSet[T](elements...)
}

// NewSetCap creates a set with a predefined capacity
func NewSetCap[T comparable](capacity int) *Set[T] {
	return WrapSet(make(map[T]struct{}, capacity))
}

// SetFromSeq creates a set with elements retrieved by the seq.
func SetFromSeq[T comparable](seq seq.Seq[T]) *Set[T] {
	if seq == nil {
		return nil
	}
	uniques := map[T]struct{}{}
	for e := range seq {
		uniques[e] = struct{}{}
	}
	return WrapSet(uniques)
}

// NewMap instantiates a map using key/value pairs
func NewMap[K comparable, V any](elements ...c.KV[K, V]) *Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// NewMapOrdered instantiates an ordered map using key/value pairs
func NewMapOrdered[K comparable, V any](elements ...c.KV[K, V]) *ordered.Map[K, V] {
	return ordered.NewMap(elements...)
}

// NewMapCap instantiates Map with a predefined capacity
func NewMapCap[K comparable, V any](capacity int) *Map[K, V] {
	return WrapMap(make(map[K]V, capacity))
}

// NewMapOf instantiates Map populated by the 'elements' map key/values
func NewMapOf[K comparable, V any](elements map[K]V) *Map[K, V] {
	return MapFromSeq2(seq2.OfMap(elements))
}

// MapFromSeq2 creates a map with elements retrieved by the seq.
func MapFromSeq2[K comparable, V any](seq seq.Seq2[K, V]) *Map[K, V] {
	if seq == nil {
		return nil
	}
	uniques := map[K]V{}
	for key, val := range seq {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// NewVector instantiates Vector populated by the 'elements' slice
func NewVector[T any](elements ...T) *Vector[T] {
	return WrapVector(slice.Clone(elements))
}

// NewVectorCap instantiates Vector with a predefined capacity
func NewVectorCap[T any](capacity int) *Vector[T] {
	return WrapVector(make([]T, 0, capacity))
}

// VectorFromSeq creates a vector with elements retrieved by the seq.
func VectorFromSeq[T any](s seq.Seq[T]) *Vector[T] {
	return WrapVector(seq.Slice(s))
}
