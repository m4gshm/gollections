// Package resolv provides values resolvers for maps thath builded by iterating over key/values loop, slice or collection
package resolv

import (
	"github.com/m4gshm/gollections/op"
)

// FirstVal - ToMap value resolver
func FirstVal[K, V any](exists bool, _ K, old, new V) V { return op.IfElse(exists, old, new) }

// LastVal - ToMap value resolver
func LastVal[K, V any](_ bool, _ K, _, new V) V { return new }

// Append - ToMap value resolver
func Append[K, V any](ok bool, k K, rv []V, v V) []V {
	return append(rv, v)
}
