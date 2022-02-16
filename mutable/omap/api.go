package omap

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/mutable/ordered"
)

//Of creates the Map with predefined elements.
func Of[K comparable, V any](elements ...*c.KV[K, V]) *ordered.Map[K, V] {
	return ordered.AsMap(elements)
}

//Empty creates the Map with zero capacity.
func Empty[K comparable, V any]() *ordered.Map[K, V] {
	return New[K, V](0)
}

//New creates the Map with a predefined capacity.
func New[K comparable, V any](capacity int) *ordered.Map[K, V] {
	return ordered.WrapMap(make([]K, 0, capacity), make(map[K]V, capacity))
}
