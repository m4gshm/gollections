package sum

import "golang.org/x/exp/constraints"

//Of returns the sum of two operands
func Of[T constraints.Ordered](a T, b T) T {
	return a + b
}
