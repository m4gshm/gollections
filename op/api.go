// Package op provides generic operations that can be used for converting or reducing collections, loops, slices
package op

import (
	"cmp"
	"fmt"

	"golang.org/x/exp/constraints"
)

// Summable is a type that supports the operator +
type Summable interface {
	cmp.Ordered | constraints.Complex | string
}

// Number is a type that supports the operators +, -, /, *
type Number interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Sum returns the sum of two operands
func Sum[T Summable](a T, b T) T {
	return a + b
}

// Sub returns the substraction of the b from the a
func Sub[T Number](a T, b T) T {
	return a - b
}

// Max returns the maximum from two operands
func Max[T cmp.Ordered](a T, b T) T {
	return IfElse(a < b, b, a)
}

// Min returns the minimum from two operands
func Min[T cmp.Ordered](a T, b T) T {
	return IfElse(a > b, b, a)
}

// IfElse returns the tru value if ok, otherwise returns the fal value
func IfElse[T any](ok bool, tru, fal T) T {
	if ok {
		return tru
	}
	return fal
}

// IfElseErrf returns the tru value if ok, otherwise returns an error creating by fmt.Errorf
func IfElseErrf[T any](ok bool, tru T, format string, a ...any) (T, error) {
	if ok {
		return tru, nil
	}
	var fal T
	return fal, fmt.Errorf(format, a...)
}

// IfElseErr returns the tru value if ok, otherwise returns the specified error
func IfElseErr[T any](ok bool, tru T, err error) (T, error) {
	if ok {
		return tru, nil
	}
	var fal T
	return fal, err
}

// IfElseGetErr returns the tru value if ok, otherwise returns an error returnet by the err function
func IfElseGetErr[T any](ok bool, tru T, err func() error) (T, error) {
	if ok {
		return tru, nil
	}
	var fal T
	return fal, err()
}

// IfElseGet returns the tru value if ok, otherwise exec the fal function and returns it result
func IfElseGet[T any](ok bool, tru T, fal func() T) T {
	if ok {
		return tru
	}
	return fal()
}

// IfGetElse executes the tru value if ok, otherwise returns the fal value
func IfGetElse[T any](ok bool, tru func() T, fal T) T {
	if ok {
		return tru()
	}
	return fal
}

// IfGetElseGet executes the tru func if ok, otherwise exec the fal function and returns it result
func IfGetElseGet[T any](ok bool, tru, fal func() T) T {
	if ok {
		return tru()
	}
	return fal()
}

// IfElseGetWithErr executes the tru func if ok, otherwise exec the fal function and returns it result
func IfElseGetWithErr[T any](ok bool, tru T, fal func() (T, error)) (T, error) {
	if ok {
		return tru, nil
	}
	return fal()
}

// IfGetElseGetErr executes the tru func if ok, otherwise exec the fal function and returns its error
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

// Compare returns -1 if o1 less than o2, 0 if equal and 1 if 01 more tha o2
func Compare[O cmp.Ordered](o1, o2 O) int {
	if o1 < o2 {
		return -1
	} else if o1 > o2 {
		return 1
	}
	return 0
}
