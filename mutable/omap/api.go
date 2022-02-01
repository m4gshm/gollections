package omap

import (
	"github.com/m4gshm/gollections/typ"
	"github.com/m4gshm/gollections/mutable/ordered"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) *ordered.Map[k, v] {
	return ordered.AsMap(elements)
}

func Empty[k comparable, v any]() *ordered.Map[k, v] {
	return New[k, v](0)
}

func New[k comparable, v any](capacity int) *ordered.Map[k, v] {
	return ordered.WrapMap(make([]k, 0, capacity), make(map[k]v, capacity))
}
