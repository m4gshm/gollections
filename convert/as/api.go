// Package as provides as.Is alias
package as

import (
	"github.com/m4gshm/gollections/convert"
	// "github.com/m4gshm/gollections/slice/iter"
)

// Is an alias of the convert.AsIs
func Is[T any](value T) T { return convert.AsIs(value) }

// Slice an alias of the convert.AsSlice
func Slice[T any](value T) []T { return convert.AsSlice(value) }

// ErrTail wraps a function of one argument and one result in a function that returns an error
func ErrTail[I, O any](f func(I) O) func(I) (O, error) {
	return func(in I) (O, error) { return f(in), nil }
}
