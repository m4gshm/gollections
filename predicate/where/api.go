// Package where provides short predicate constructors
package where

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
)

// Match - where.Match alias for the predicate.Match
func Match[From, To any](getter func(From) To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Match(getter, condition)
}

// Any - where.Any alias for the predicate.MatchAny
func Any[From, To any](getter func(From) []To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.MatchAny(getter, condition)
}

// Eq creates predicate thet checks equality of a strcut property value to the specified example
func Eq[From any, To comparable](getter func(From) To, example To) predicate.Predicate[From] {
	return Match(getter, eq.To(example))
}

// Not creates negate condition for a  a strcut property
func Not[From, To any](getter func(From) To, condition predicate.Predicate[To]) predicate.Predicate[From] {
	return predicate.Not(Match(getter, condition))
}

func Key[K comparable, V any, M ~map[K]V](key K) predicate.Predicate[M] {
	return func(m M) bool { _, ok := m[key]; return ok }
}
