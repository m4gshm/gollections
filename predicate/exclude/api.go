package exclude

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/one"
)

// All - exclude.All creates a predicate that tests  if a value is not in the excluded values
func All[T comparable](excluded ...T) predicate.Predicate[T] {
	return predicate.Not(one.Of(excluded...))
}
