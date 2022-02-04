//Package omap provides the ordered map container implementation
package omap

import (
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/map_"
)

func Of[k comparable, v any](elements ...*map_.KV[k, v]) *ordered.Map[k, v] {
	return ordered.ConvertKVsToMap(elements)
}

func New[k comparable, v any](elements map[k]v) *ordered.Map[k, v] {
	return ordered.NewMap(elements)
}
