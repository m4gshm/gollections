// Package group provides short aliases for functions that are used to group key/values retieved from a source
package group

import (
	"github.com/m4gshm/gollections/map_"
)

// OfLoop - group.OfLoop synonym for the map_.GroupOfLoop.
func OfLoop[S any, K comparable, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) (map[K][]V, error) {
	return map_.GroupOfLoop(source, hasNext, getNext)
}
