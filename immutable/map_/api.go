// Package map_ provides immutable.Map constructors
package map_ //nilint

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/kv/loop"
)

// Of instantiates a ap from the specified key/value pairs
func Of[K comparable, V any](elements ...c.KV[K, V]) immutable.Map[K, V] {
	return immutable.ConvertKVsToMap(elements)
}

// New instantiates Map and copies elements to it
func New[K comparable, V any](elements map[K]V) immutable.Map[K, V] {
	return immutable.NewMap(elements)
}

// From instantiates a map with key/values retrieved by the 'next' function.
// The next returns a key/value pairs with true or zero values with false if there are no more elements.
func From[K comparable, V any](next func() (K, V, bool)) immutable.Map[K, V] {
	return immutable.WrapMap(loop.ToMap(next))
}
