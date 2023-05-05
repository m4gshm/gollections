// Package ordered provides mutable ordered collection implementations
package ordered

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/clone"
)

// NewSet instantiates set and copies elements to it
func NewSet[T comparable](elements ...T) *Set[T] {
	var (
		l       = len(elements)
		uniques = make(map[T]int, l)
		order   = make([]T, 0, l)
	)
	pos := 0
	for _, e := range elements {
		order, pos = addToSet(e, uniques, order, pos)
	}
	return WrapSet(order, uniques)
}

// NewSetCap creates a set with a predefined capacity
func NewSetCap[T comparable](capacity int) *Set[T] {
	return WrapSet(make([]T, 0, capacity), make(map[T]int, capacity))
}

// SetFromLoop creates a set with elements retrieved by the 'next' function.
// The next returns an element with true or zero value with false if there are no more elements.
func SetFromLoop[T comparable](next func() (T, bool)) *Set[T] {
	var (
		uniques = map[T]int{}
		order   []T
		pos     = 0
	)
	for e, ok := next(); ok; e, ok = next() {
		order, pos = addToSet(e, uniques, order, pos)
	}
	return WrapSet(order, uniques)
}

// NewMap instantiates a map using key/value pairs
func NewMap[K comparable, V any](elements ...c.KV[K, V]) *Map[K, V] {
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

// NewMapCap instantiates Map with a predefined capacity
func NewMapCap[K comparable, V any](capacity int) *Map[K, V] {
	return WrapMap(make([]K, capacity), make(map[K]V, capacity))
}

// NewMapOf instantiates Map populated by the 'elements' map key/values
func NewMapOf[K comparable, V any](order []K, elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, key := range order {
		uniques[key] = elements[key]
	}
	return WrapMap(clone.Of(order), uniques)
}

// MapFromLoop creates a map with elements retrieved converter the 'next' function
func MapFromLoop[K comparable, V any](next func() (K, V, bool)) *Map[K, V] {
	var (
		uniques = map[K]V{}
		order   = []K{}
	)
	for key, val, ok := next(); ok; key, val, ok = next() {
		order = addToMap(key, val, order, uniques)
	}
	return WrapMap(order, uniques)
}
