// Package op provides generic operations that can be used for converting or reducing collections, loops, slices
package op

import (
	"cmp"

	"github.com/m4gshm/gollections/op"
)

// Sum returns the sum of two operands.
func Sum[T op.Summable](a T, b T) (T, error) {
	return op.Sum(a, b), nil
}

// Sub returns the subtraction of the b from the a.
func Sub[T op.Number](a T, b T) (T, error) {
	return op.Sub(a, b), nil
}

// Max returns the maximum from two operands.
func Max[T cmp.Ordered](a T, b T) (T, error) {
	return IfElse(a < b, b, a)
}

// Min returns the minimum from two operands.
func Min[T cmp.Ordered](a T, b T) (T, error) {
	return IfElse(a > b, b, a)
}

// IfElse returns the tru value if ok, otherwise return the fal value.
func IfElse[T any](ok bool, tru, fal T) (T, error) {
	if ok {
		return tru, nil
	}
	return fal, nil
}

// IfGetElseGet executes the tru func if ok, otherwise exec the fal function and returns it result.
func IfGetElseGet[T any](ok bool, tru, fal func() (T, error)) (T, error) {
	if ok {
		return tru()
	}
	return fal()
}
