package check

import (
	"reflect"
)

// Nil checks a reference for nil value.
func Nil[T any](val *T) bool {
	return val == nil
}

// NotNil checks a reference for no nil value.
func NotNil[T any](val *T) bool {
	return !Nil(val)
}

// Zero checks that a value is zero.
func Zero[T any](val T) bool {
	return reflect.ValueOf(val).IsZero()
}

// Empty checks that a slice is empty.
func Empty[T Slice | []any](val T) bool {
	return len(val) == 0
}

// NotEmpty checks that a slice is not empty.
func NotEmpty[C []T, T any](val C) bool {
	return len(val) > 0
}

// EmptyMap checks that a slice is ampty.
func EmptyMap[K comparable, V any](val map[K]V) bool {
	return len(val) == 0
}

// NotEmptyMap checks that a slice is not empty.
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
