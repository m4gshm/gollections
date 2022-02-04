//Package map_ provides the unordered map container implementation
package map_ //nilint

import (
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/map_"
)

func Of[k comparable, v any](elements ...*map_.KV[k, v]) *immutable.Map[k, v] {
	return immutable.ConvertKVsToMap(elements)
}

func New[k comparable, v any](elements map[k]v) *immutable.Map[k, v] {
	return immutable.NewMap(elements)
}
