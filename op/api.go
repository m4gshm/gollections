// Package op provides generic operations that can be used for converting or reducing collections, loops, slices
package op

import (
	"github.com/m4gshm/gollections/c"
	"golang.org/x/exp/constraints"
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

// Empty checks the val is empty
func Empty[T Slice | []any](val T) bool {
	return len(val) == 0
}

// NotEmpty checks the val is not empty
func NotEmpty[C []T, T any](val C) bool {
	return len(val) > 0
}

// EmptyMap checks the val map is empty
func EmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) == 0
}

// NotEmptyMap checks the val map is not empty
func NotEmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) > 0
}

// Slice is the constraint included all slice types
type Slice interface {
	~[]any | ~[]uintptr |
		~[]int | ~[]int8 | []int16 | []int32 | []int64 |
		~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 |
		~[]float32 | ~[]float64 |
		~[]complex64 | ~[]complex128 |
		~[]string | ~string
}
