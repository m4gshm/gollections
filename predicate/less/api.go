package less

import (
	"github.com/m4gshm/gollections/predicate"
	"golang.org/x/exp/constraints"
)

// Than - less.Than creates a predicate that can be used to test if a value is less than the expected
func Than[T constraints.Ordered](min T) predicate.Predicate[T] {
	return func(v T) bool { return v < min }
}

// OrEq - less.OrEq creates a predicate that can be used to test if a value is less than or equal to the expected
func OrEq[T constraints.Ordered](min T) predicate.Predicate[T] {
	return func(v T) bool { return v <= min }
}
