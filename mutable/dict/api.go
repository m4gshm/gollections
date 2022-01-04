package dict

import (
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) mutable.Map[k, v, typ.Iterator[*typ.KV[k,v]]] {
	return ToOrderedMap(elements)
}

func New[k comparable, v any](capacity int) mutable.Map[k, v, typ.Iterator[*typ.KV[k,v]]] {
	return NewOrderedMap[k, v](capacity)
}