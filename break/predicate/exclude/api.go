// Package exclude provides predicate builders
package exclude

import (
	"github.com/m4gshm/gollections/break/predicate"
	"github.com/m4gshm/gollections/break/predicate/one"
)

// All - exclude.All creates a predicate that tests  if a value is not in the excluded values
func All[T comparable](excluded ...T) predicate.Predicate[T] {
	return predicate.Not(one.Of(excluded...))
}
