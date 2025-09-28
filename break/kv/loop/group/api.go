// Package group provides short aliases for functions thath are used to group key/value pairs retrieved by a seq
package group

import (
	"github.com/m4gshm/gollections/break/kv/loop"
)

// Of is a short alias for loop.Group
func Of[K comparable, V any](next func() (K, V, bool, error)) (map[K][]V, error) {
	return loop.Group(next)
}
