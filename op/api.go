package op

import (
	"constraints"
)

type Binary[T any] func(T, T) T

type Quaternary[t1, t2 any] func(t1, t2, t1, t2) (t1, t2)

func Sum[T constraints.Ordered](a T, b T) T {
	return a + b
}
