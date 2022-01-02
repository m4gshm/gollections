package omap

import (
	"github.com/m4gshm/container/immutable"
	"github.com/m4gshm/container/typ"
)

func Of[k comparable, v any](values ...*typ.KV[k, v]) typ.Map[k, v] {
	return immutable.NewOrderedMap(values)
}

func New[k comparable, v any](values []*typ.KV[k, v]) typ.Map[k, v] {
	return immutable.NewOrderedMap(values)
}