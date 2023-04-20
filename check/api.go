// Package check provides common predicate functions
package check

import (
	"reflect"
)

// Nil checks whether the reference is nil
func Nil[T any](val *T) bool {
	return val == nil
}

// NotNil checks whether the reference is not nil
func NotNil[T any](val *T) bool {
	return !Nil(val)
}

// Zero checks whether the value is zero
func Zero[T any](value T) bool {
	return reflect.ValueOf(value).IsZero()
}

// Empty checks whether the slice or string is empty
func Empty[S Slice | []any](s S) bool {
	return len(s) == 0
}

// NotEmpty checks whether the slice is no empty
func NotEmpty[S Slice | []any](slice S) bool {
	return len(slice) > 0
}

// EmptyMap checks whether the map is empty
func EmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) == 0
}

// NotEmptyMap checks whether the map is not empty
func NotEmptyMap[K comparable, V any](m map[K]V) bool {
	return len(m) > 0
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
