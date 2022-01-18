package map_

import (
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) mutable.Map[k, v, typ.KVIterator[k, v]] {
	return ToOrderedMap(elements)
}

func New[k comparable, v any](capacity int) mutable.Map[k, v, typ.KVIterator[k, v]] {
	return NewOrderedMap[k, v](capacity)
}
