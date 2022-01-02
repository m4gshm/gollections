package dict

import (
	"github.com/m4gshm/container/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) typ.Map[k, v] {
	return NewOrderedMap(elements)
}

func New[k comparable, v any](elements []*typ.KV[k, v]) typ.Map[k, v] {
	return NewOrderedMap(elements)
}