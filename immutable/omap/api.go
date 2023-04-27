// Package omap provides immutale ordered.Map constructors
package omap

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/kv/loop"
)

// Of instantiates a ap from the specified key/value pairs
func Of[K comparable, V any](pairs ...c.KV[K, V]) ordered.Map[K, V] {
	return ordered.NewMapKV(pairs)
}

// New instantiates a map and copies the elements to it
func New[K comparable, V any](elements map[K]V, order []K) ordered.Map[K, V] {
	return ordered.NewMap(order, elements)
}

// From instantiates a map with key/values retrieved by the 'next' function.
// The next returns a key/value pairs with true or zero values with false if there are no more elements.
func From[K comparable, V any](next func() (K, V, bool)) ordered.Map[K, V] {
	return ordered.NewMapKV(loop.ToSlice(next))
}
