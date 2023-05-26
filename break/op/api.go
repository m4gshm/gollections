// Package op provides generic operations that can be used for converting or reducing collections, loops, slices
package op

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/op"
)

// Sum returns the sum of two operands
func Sum[T c.Summable](a T, b T) (T, error) {
	return op.Sum(a, b), nil
}

// Sub returns the substraction of the b from the a
func Sub[T c.Number](a T, b T) (T, error) {
	return op.Sub(a, b), nil
}

// Max returns the maximum from two operands
func Max[T constraints.Ordered](a T, b T) (T, error) {
	return IfElse(a < b, b, a)
}

// Min returns the minimum from two operands
func Min[T constraints.Ordered](a T, b T) (T, error) {
	return IfElse(a > b, b, a)
}

// IfElse returns the tru value if ok, otherwise return the fal value
func IfElse[T any](ok bool, tru, fal T) (T, error) {
	if ok {
		return tru, nil
	}
	return fal, nil
}

// IfDoElse exececutes the tru func if ok, otherwise exec the fal function and returns it result
func IfDoElse[T any](ok bool, tru, fal func() (T, error)) (T, error) {
	if ok {
		return tru()
	}
	return fal()
}
