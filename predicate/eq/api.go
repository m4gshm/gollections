package eq

import "github.com/m4gshm/gollections/predicate"

// Eq makes a predicate to test for equality.
func To[T comparable](v T) predicate.Predicate[T] {
	return predicate.Eq(v)
}