// Package match provides short predicate constructors
package match

import "github.com/m4gshm/gollections/predicate"

// To - match.To alias for the predicate.Match
func To[From, To any](getter func(From) To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Match(getter, condition)
}

// Any - match.Any alias for the predicate.MatchAny
func Any[From, To any](getter func(From) []To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.MatchAny(getter, condition)
}
