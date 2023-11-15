// Package resolv provides values resolvers for maps that builded by ToMap-converter functions
package resolv

import (
	"cmp"
	"slices"

	"github.com/m4gshm/gollections/op"
)

// First keeps the first value of a key
func First[K, V any](exists bool, _ K, old, new V) V { return op.IfElse(exists, old, new) }

// Last retrieves the last value of a key
func Last[K, V any](_ bool, _ K, _, new V) V { return new }

// Slice puts the values of one key into a slice
func Slice[K, V any](_ bool, _ K, rv []V, v V) []V {
	return append(rv, v)
}

// SortedSlice puts the values of one key into a sorted slice
func SortedSlice[K, V cmp.Ordered](_ bool, _ K, rv []V, v V) []V {
	i, _ := slices.BinarySearch[[]V, V](rv, v)
	r := append(append(rv[:i], v), rv[i:]...)
	return r
}
