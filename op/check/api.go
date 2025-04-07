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

func Empty[TS ~[]T, T any](elements TS) bool {
	return len(elements) == 0
}

func NotEmpty[TS ~[]T, T any](elements TS) bool {
	return !Empty(elements)
}

func EmptyStr(s string) bool {
	return string_.Empty(s)
}

func NotEmptyStr(s string) bool {
	return !EmptyStr(s)
}
