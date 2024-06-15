package is

import (
	"github.com/m4gshm/gollections/predicate"
)

// Not negates the 'p' predicate
func Not[T any](p predicate.Predicate[T]) predicate.Predicate[T] {
	return predicate.Not(p)
}
