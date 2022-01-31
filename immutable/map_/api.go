//map_ package oset provides the unordered map container implementation
package map_

import (
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) *immutable.Map[k, v] {
	return immutable.ConvertKVsToMap(elements)
}

func New[k comparable, v any](elements map[k]v) *immutable.Map[k, v] {
	return immutable.NewMap(elements)
}
