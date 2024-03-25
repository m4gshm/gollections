// Package filter provides helpers for filtering keys or values of a map
package filter

import "github.com/m4gshm/gollections/map_"

// Keys an alias of the map_.FilterKeys
func Keys[M ~map[K]V, K comparable, V any](m M, filter func(K) bool) map[K]V {
	return map_.FilterKeys(m, filter)
}

// Values an alias of the map_.FilterValues
func Values[M ~map[K]V, K comparable, V any](m M, filter func(V) bool) map[K]V {
	return map_.FilterValues(m, filter)
}
