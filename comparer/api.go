package comparer

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/op"
)

func Of[T any, O constraints.Ordered](converter func(T) O) func(T, T) int {
	return func(e1, e2 T) int {
		return op.Compare(converter(e1), converter(e2))
	}
}

func Reverse[T any, O constraints.Ordered](converter func(T) O) func(T, T) int {
	return func(e1, e2 T) int {
		return -op.Compare(converter(e1), converter(e2))
	}
}
