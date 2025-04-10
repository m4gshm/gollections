// Package check provides common predicate functions
package check

import (
	"reflect"

	"github.com/m4gshm/gollections/op/string_"
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

// Empty checks that the slice does not contain elements
func Empty[TS ~[]T, T any](slice TS) bool {
	return len(slice) == 0
}

// NotEmpty checks that the slice contains elements
func NotEmpty[TS ~[]T, T any](slice TS) bool {
	return !Empty(slice)
}

// EmptyStr checks whether the specified string is empty
func EmptyStr(s string) bool {
	return string_.Empty(s)
}

// NotEmptyStr checks whether the specified string is not empty
func NotEmptyStr(s string) bool {
	return !EmptyStr(s)
}

// EmptyMap checks if the elements map is empty
func EmptyMap[M ~map[K]V, K comparable, V any](elements M) bool {
	return len(elements) == 0
}

// NotEmptyMap checks if the elements map is not empty
func NotEmptyMap[M ~map[K]V, K comparable, V any](elements M) bool {
	return !EmptyMap(elements)
}
