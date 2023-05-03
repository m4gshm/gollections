// Package more provides predicate builders
package more

import (
	"github.com/m4gshm/gollections/break/predicate"
	"golang.org/x/exp/constraints"
)

// Than - more.Than creates a predicate that can be used to test if a value is greater than the expected
func Than[T constraints.Ordered](expected T) predicate.Predicate[T] {
	return func(v T) (bool, error) { return v > expected, nil }
}

// OrEq - more.OrEq creates a predicate that can be used to test if a value is greater than or equal to the expected
func OrEq[T constraints.Ordered](expected T) predicate.Predicate[T] {
	return func(v T) (bool, error) { return v >= expected, nil }
}
