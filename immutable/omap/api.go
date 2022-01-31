//package omap provides the ordered map container implementation
package omap

import (
	"github.com/m4gshm/gollections/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) *OrderedMap[k, v] {
	return Convert(elements)
}

func New[k comparable, v any](elements map[k]v) *OrderedMap[k, v] {
	return ConvertMap(elements)
}
