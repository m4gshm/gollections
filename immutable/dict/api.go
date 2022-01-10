package dict

import (
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) immutable.Map[k, v] {
	return NewOrderedMap(elements)
}

func New[k comparable, v any](elements []*typ.KV[k, v]) immutable.Map[k, v] {
	return NewOrderedMap(elements)
}
