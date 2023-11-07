// Package filter provides helpers for filtering keys or values of a map
package filter

import "github.com/m4gshm/gollections/map_"

func Values[M ~map[K]V, K comparable, V any](m M, filter func(V) bool) map[K]V {
	return map_.FilterValues(m, filter)
}

// Key adapts a key appliable predicate to a key\value one
func Key[V, K any](predicate func(K) bool) (out func(K, V) bool) {
	if predicate != nil {
		out = func(key K, val V) bool { return predicate(key) }
	}
	return
}

// Value adapts a value appliable predicate to a key\value one
func Value[K, V any](predicate func(V) bool) (out func(K, V) bool) {
	if predicate != nil {
		out = func(key K, val V) bool { return predicate(val) }
	}
	return
}
