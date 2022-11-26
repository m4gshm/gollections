package map_

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/mutable"
)

// Of instantiates Map with predefined elements.
func Of[K comparable, V any](elements ...c.KV[K, V]) mutable.Map[K, V] {
	return mutable.AsMap(elements)
}

// Empty instantiates Map with zero capacity.
func Empty[K comparable, V any]() mutable.Map[K, V] {
	return New[K, V](0)
}

// New instantiates Map with a predefined capacity.
func New[K comparable, V any](capacity int) mutable.Map[K, V] {
	return mutable.NewMap[K, V](0)
}
