package op

import "constraints"

type Binary[T any] func(T, T) T

func Sum[T constraints.Ordered](a T, b T) T {
	return a + b
}
