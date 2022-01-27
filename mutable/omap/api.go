package omap

import (
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/typ"
)

func Of[k comparable, v any](elements ...*typ.KV[k, v]) mutable.Map[k, v] {
	return Convert(elements)
}

func Empty[k comparable, v any]() mutable.Map[k, v] {
	return New[k, v](0)
}

func New[k comparable, v any](capacity int) mutable.Map[k, v] {
	return Wrap(make([]*k, 0, capacity), make(map[k]v, capacity))
}