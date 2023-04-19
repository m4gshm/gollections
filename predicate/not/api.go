package not

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/eq"
)

// not.Eq makes a reverse Eq predicate.
func Eq[T comparable](v T) predicate.Predicate[T] {
	return predicate.Not(eq.To(v))
}
