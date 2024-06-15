package is

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
)

// Not negates the 'p' predicate
func Not[T any](p predicate.Predicate[T]) predicate.Predicate[T] {
	return predicate.Not(p)
}

// Eq creates a predicate to test for equality
func Eq[T comparable](v T) predicate.Predicate[T] {
	return eq.To(v)
}
