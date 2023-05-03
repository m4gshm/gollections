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

func Empty[From, To any](flattener func(From) []To) predicate.Predicate[From] {
	return predicate.Not(predicate.Empty(flattener))
}

// func Empty2[ T string ](flattener func(From) []To) predicate.Predicate[From] {
// 	return predicate.Not(predicate.Empty(flattener))
// }
