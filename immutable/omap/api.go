//Package omap provides the ordered map container implementation
package omap

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
)

//Of creates the Map with predefined elements.
func Of[K comparable, V any](elements ...*c.KV[K, V]) *ordered.Map[K, V] {
	return ordered.ConvertKVsToMap(elements)
}

//New creates the Map and copies elements to it.
func New[K comparable, V any](elements map[K]V) *ordered.Map[K, V] {
	return ordered.NewMap(elements)
}
