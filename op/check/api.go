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
