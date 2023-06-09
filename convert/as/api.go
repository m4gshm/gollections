// Package as provides as.Is alias
package as

import (
	"github.com/m4gshm/gollections/convert"
)

// Is an alias of the convert.AsIs
func Is[T any](value T) T { return convert.AsIs(value) }

// Slice an alias of the convert.AsSlice
func Slice[T any](value T) []T { return convert.AsSlice(value) }

// ErrTail wraps a function of one argument and one result in a function that returns an error
func ErrTail[I, O any](f func(I) O) func(I) (O, error) {
	return func(in I) (O, error) { return f(in), nil }
}

// Ptr converts a value to the value pointer
func Ptr[T any](value T) *T { return convert.Ptr(value) }

// Val returns a value referenced by the pointer or the zero value if the pointer is nil
func Val[T any](pointer *T) T { return convert.PtrVal(pointer) }
