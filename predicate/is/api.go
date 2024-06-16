package is

import (
	"github.com/m4gshm/gollections/predicate"
)

// Not - is.Not creates a 'not p' predicate.
func Not[T any](p predicate.Predicate[T]) predicate.Predicate[T] {
	return predicate.Not(p)
}
