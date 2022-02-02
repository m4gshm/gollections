package omap

import (
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable/ordered"
)

func Of[k comparable, v any](elements ...*map_.KV[k, v]) *ordered.Map[k, v] {
	return ordered.AsMap(elements)
}

func Empty[k comparable, v any]() *ordered.Map[k, v] {
	return New[k, v](0)
}

func New[k comparable, v any](capacity int) *ordered.Map[k, v] {
	return ordered.WrapMap(make([]k, 0, capacity), make(map[k]v, capacity))
}
