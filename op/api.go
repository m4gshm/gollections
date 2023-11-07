// Package op provides generic operations that can be used for converting or reducing collections, loops, slices
package op

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/c"
)

// Sum returns the sum of two operands
func Sum[T c.Summable](a T, b T) T {
	return a + b
}

// Sub returns the substraction of the b from the a
func Sub[T c.Number](a T, b T) T {
	return a - b
}

// Max returns the maximum from two operands
func Max[T constraints.Ordered](a T, b T) T {
	return IfElse(a < b, b, a)
}

// Min returns the minimum from two operands
func Min[T constraints.Ordered](a T, b T) T {
	return IfElse(a > b, b, a)
}

// IfElse returns the tru value if ok, otherwise return the fal value
func IfElse[T any](ok bool, tru, fal T) T {
	if ok {
		return tru
	}
	return fal
}

// IfElseErr returns the tru value if ok, otherwise return the specified error
func IfElseErr[T any](ok bool, tru T, err error) (T, error) {
	if ok {
		return tru, nil
	}
	var fal T
	return fal, err
}

// IfGetElse exececutes the tru func if ok, otherwise exec the fal function and returns it result
func IfGetElse[T any](ok bool, tru, fal func() T) T {
	if ok {
		return tru()
	}
	return fal()
}

// IfGetElseGetErr exececutes the tru func if ok, otherwise exec the fal function and returns its error
func IfGetElseGetErr[T any](ok bool, tru func() T, fal func() error) (T, error) {
	if ok {
		return tru(), nil
	}
	var no T
	return no, fal()
}

// Get calls the getter and returns the result
func Get[T any](getter func() T) T {
	return getter()
}

func Compare[O constraints.Ordered](o1, o2 O) int {
	if o1 < o2 {
		return -1
	} else if o1 > o2 {
		return 1
	}
	return 0
}
