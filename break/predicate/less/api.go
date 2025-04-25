// Package less provides predicate builders
package less

import (
	"github.com/m4gshm/gollections/break/predicate"
	"golang.org/x/exp/constraints"
)

// Than - less.Than creates a predicate that can be used to test if a value is less than the expected
func Than[T constraints.Ordered](threshold T) predicate.Predicate[T] {
	return func(v T) (bool, error) { return v < threshold, nil }
}

// OrEq - less.OrEq creates a predicate that can be used to test if a value is less than or equal to the expected
func OrEq[T constraints.Ordered](threshold T) predicate.Predicate[T] {
	return func(v T) (bool, error) { return v <= threshold, nil }
}
