// Package group provides short aliases for functions that are used to group collection elements
package group

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

// Of converts the 'elements' slice of key\value pairs.
func Of[TS ~[]c.KV[K, V], K comparable, V any](elements TS) map[K][]V {
	if elements == nil {
		return nil
	}
	return slice.Group(elements, c.KV[K, V].Key, c.KV[K, V].Value)
}
