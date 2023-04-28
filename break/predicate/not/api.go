// Package not provides negalive predicate builders like 'not equals to'
package not

import (
	"github.com/m4gshm/gollections/break/predicate"
	"github.com/m4gshm/gollections/break/predicate/eq"
)

// Eq - not.Eq makes reverse of the eq.To predicate
func Eq[T comparable](v T) predicate.Predicate[T] {
	return predicate.Not(eq.To(v))
}
