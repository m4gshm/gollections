// Package map_ provides immutale ordered.Map constructors
package map_

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
)

// Of instantiates a ap from the specified key/value pairs
func Of[K comparable, V any](pairs ...c.KV[K, V]) ordered.Map[K, V] {
	return ordered.NewMap(pairs...)
}

// New instantiates a map and copies the elements to it
func New[K comparable, V any](elements map[K]V, order []K) ordered.Map[K, V] {
	return ordered.NewMapOf(order, elements)
}

// From instantiates a map with key/values retrieved by the 'next' function.
// The next returns a key/value pairs with true or zero values with false if there are no more elements.
func From[K comparable, V any](next func() (K, V, bool)) ordered.Map[K, V] {
	return ordered.MapFromLoop(next)
}

// Resolv collects key\value elements to an ordered.Map by iterating over the elements with resolving of duplicated key values
// TODO: needs test
func Resolv[TS ~[]T, T any, K comparable, V, VR any](elements TS, keyExtractor func(T) K, valExtractor func(T) V, resolver func(bool, K, VR, V) VR) ordered.Map[K, VR] {
	l := len(elements)
	m := make(map[K]VR)
	order := make([]K, 0, l)
	for _, e := range elements {
		k, v := keyExtractor(e), valExtractor(e)
		exists, ok := m[k]
		if !ok {
			order = append(order, k)
		}
		m[k] = resolver(ok, k, exists, v)
	}
	return ordered.WrapMap(order, m)
}
