package op

type Binary[T any] func(T, T) T

type Additionable interface {
	string | int | int8 | int32 | int64 | uint | uint8 | uint32 | uint64 | float32 | float64
}

func Sum[T Additionable](a T, b T) T {
	return a + b
}
