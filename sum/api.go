package sum

import "golang.org/x/exp/constraints"

//Of returns two elements sum
func Of[T constraints.Ordered](a T, b T) T {
	return a + b
}
