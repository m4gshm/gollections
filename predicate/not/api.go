// Package not provides negalive predicate builders like 'not equals to'
package not

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
)

// Eq - not.Eq makes reverse of the eq.To predicate
func Eq[T comparable](v T) predicate.Predicate[T] {
	return predicate.Not(eq.To(v))
}

// Match - not.Match alias of predicate.Not
func Match[From, To any](getter func(From) To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Not(predicate.Match(getter, condition))
}
