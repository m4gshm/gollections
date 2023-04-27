package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

func Of[TS ~[]c.KV[K, V], K comparable, V any](elements TS) map[K][]V {
	if elements == nil {
		return nil
	}
	return slice.Group(elements, c.KV[K, V].Key, c.KV[K, V].Value)
}
