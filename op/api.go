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
	if a < b {
		return b
	}
	return a
}

// Min returns the minimum from two operands
func Min[T constraints.Ordered](a T, b T) T {
	if a > b {
		return b
	}
	return a
}
