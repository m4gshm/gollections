//Package map_ provides the unordered map container implementation
package map_ //nilint

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
)

//Of creates the Map with predefined elements.
func Of[k comparable, v any](elements ...c.KV[k, v]) immutable.Map[k, v] {
	return immutable.ConvertKVsToMap(elements)
}

//New creates the Map and copies elements to it.
func New[k comparable, v any](elements map[k]v) immutable.Map[k, v] {
	return immutable.NewMap(elements)
}
