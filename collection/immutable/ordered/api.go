// Package ordered provides immutable ordered collection implementations
package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/slice/clone"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) Set[T] {
	return SetFromSeq(seq.Of(elements...))
}

// SetFromSeq creates a set with elements retrieved by the seq.
func SetFromSeq[T comparable](seq seq.Seq[T]) Set[T] {
	if seq == nil {
		return Set[T]{}
	}
	var (
		uniques = map[T]struct{}{}
		order   []T
	)

	for e := range seq {
		order = addToSet(e, uniques, order)
	}
	return WrapSet(order, uniques)
}

// NewMap instantiates a map using key/value pairs
func NewMap[K comparable, V any](elements ...c.KV[K, V]) Map[K, V] {
	var (
		l       = len(elements)
		uniques = make(map[K]V, l)
		order   = make([]K, 0, l)
	)
	for _, kv := range elements {
		order = addToMap(kv.Key(), kv.Value(), order, uniques)
	}
	return WrapMap(order, uniques)
}

// NewMapOf instantiates Map populated by the 'elements' map key/values
func NewMapOf[K comparable, V any](order []K, elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, key := range order {
		uniques[key] = elements[key]
	}
	return WrapMap(clone.Of(order), uniques)
}

// MapFromSeq2 creates a map with elements retrieved by the seq.
func MapFromSeq2[K comparable, V any](seq seq.Seq2[K, V]) Map[K, V] {
	if seq == nil {
		return Map[K, V]{}
	}
	var (
		uniques = map[K]V{}
		order   = []K{}
	)
	for key, val := range seq {
		order = addToMap(key, val, order, uniques)
	}
	return WrapMap(order, uniques)
}
