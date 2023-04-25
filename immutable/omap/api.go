// Package omap provides immutale ordered.Map constructors
package omap

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
)

// Of instantiates Map with predefined elements.
func Of[K comparable, V any](elements ...c.KV[K, V]) ordered.Map[K, V] {
	return ordered.NewMapKV(elements)
}

// New instantiates Map and copies elements to it.
func New[K comparable, V any](elements map[K]V, order []K) ordered.Map[K, V] {
	return ordered.NewMap(elements, order)
}
