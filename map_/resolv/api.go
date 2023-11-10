// Package resolv provides values resolvers for maps thath builded by iterating over key/values loop, slice or collection
package resolv

import (
	"cmp"
	"slices"

	"github.com/m4gshm/gollections/op"
)

// First - ToMap value resolver
func First[K, V any](exists bool, _ K, old, new V) V { return op.IfElse(exists, old, new) }

// Last - ToMap value resolver
func Last[K, V any](_ bool, _ K, _, new V) V { return new }

// Slice - ToMap value resolver
func Slice[K, V any](_ bool, _ K, rv []V, v V) []V {
	return append(rv, v)
}

func SortedSlice[K, V cmp.Ordered](_ bool, _ K, rv []V, v V) []V {
	i, _ := slices.BinarySearch[[]V, V](rv, v)
	r := append(append(rv[:i], v), rv[i:]...)
	return r
}
