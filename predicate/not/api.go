// Package not provides negalive predicate builders like 'not equals to'
package not

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
	"github.com/m4gshm/gollections/predicate/is"
)

// Eq - not.Eq makes a 'not eq.To' predicate
func Eq[T comparable](v T) predicate.Predicate[T] {
	return is.Not(eq.To(v))
}

// Match - not.Match makes a 'not predicate.Match' predicate
func Match[From, To any](getter func(From) To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return is.Not(predicate.Match(getter, condition))
}
