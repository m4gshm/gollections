// Package group provides short aliases for functions thath are used to group key/value pairs retrieved by a loop
package group

import (
	"github.com/m4gshm/gollections/kv/loop"
)

// Of is a short alias for loop.Group
func Of[K comparable, V any](next func() (K, V, bool)) map[K][]V {
	return loop.Group(next)
}
