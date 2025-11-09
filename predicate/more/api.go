// Package more provides predicate builders
package more

import (
	"cmp"

	"github.com/m4gshm/gollections/predicate"
)

// Than - more.Than creates a predicate that can be used to test if a value is greater than the expected
func Than[T cmp.Ordered](expected T) predicate.Predicate[T] {
	return func(v T) bool { return v > expected }
}

// OrEq - more.OrEq creates a predicate that can be used to test if a value is greater than or equal to the expected
func OrEq[T cmp.Ordered](expected T) predicate.Predicate[T] {
	return func(v T) bool { return v >= expected }
}
