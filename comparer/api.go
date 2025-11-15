// Package comparer provides builders of slices.CompareFunc comparison functions
package comparer

import (
	"cmp"

	"github.com/m4gshm/gollections/op"
)

// Of creates a comparer for orderable values obtained by the converter
func Of[T any, O cmp.Ordered](converter func(T) O) func(T, T) int {
	return func(e1, e2 T) int {
		return op.Compare(converter(e1), converter(e2))
	}
}

// Reverse creates a descending comparer for orderable values obtained by the converter
func Reverse[T any, O cmp.Ordered](converter func(T) O) func(T, T) int {
	return func(e1, e2 T) int {
		return -op.Compare(converter(e1), converter(e2))
	}
}
