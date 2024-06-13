// Package ordered provides immutable ordered collection implementations
package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/clone"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) Set[T] {
	var (
		l       = len(elements)
		uniques = make(map[T]struct{}, l)
		order   = make([]T, 0, l)
	)
	for _, e := range elements {
		order = addToSet(e, uniques, order)
	}
	return WrapSet(order, uniques)
}

// SetFromLoop creates a set with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func SetFromLoop[T comparable](next func() (T, bool)) Set[T] {
	if next == nil {
		return Set[T]{}
	}
	var (
		uniques = map[T]struct{}{}
		order   []T
	)
	for e, ok := next(); ok; e, ok = next() {
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

// MapFromLoop creates a map with elements retrieved converter the 'next' function
func MapFromLoop[K comparable, V any](next func() (K, V, bool)) Map[K, V] {
	if next == nil {
		return Map[K, V]{}
	}
	var (
		uniques = map[K]V{}
		order   = []K{}
	)
	for key, val, ok := next(); ok; key, val, ok = next() {
		order = addToMap(key, val, order, uniques)
	}
	return WrapMap(order, uniques)
}
