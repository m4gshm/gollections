package not

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
)

// Eq - not.Eq makes reverse of the eq.To predicate
func Eq[T comparable](v T) predicate.Predicate[T] {
	return predicate.Not(eq.To(v))
}
