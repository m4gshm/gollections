// Package eq provides predicate builder short aliases
package eq

import "github.com/m4gshm/gollections/predicate"

// To - eq.To creates a predicate to test for equality to the 'v'
func To[T comparable](v T) predicate.Predicate[T] {
	return predicate.Eq(v)
}
