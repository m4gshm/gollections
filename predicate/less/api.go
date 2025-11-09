// Package less provides predicate builders
package less

import (
	"cmp"

	"github.com/m4gshm/gollections/predicate"
)

// Than - less.Than creates a predicate that can be used to test if a value is less than the expected
func Than[T cmp.Ordered](threshold T) predicate.Predicate[T] {
	return func(v T) bool { return v < threshold }
}

// OrEq - less.OrEq creates a predicate that can be used to test if a value is less than or equal to the expected
func OrEq[T cmp.Ordered](threshold T) predicate.Predicate[T] {
	return func(v T) bool { return v <= threshold }
}
