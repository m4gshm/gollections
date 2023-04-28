// Package one provides predicate builders
package one

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/break/predicate"
	"github.com/m4gshm/gollections/immutable/set"
)

// Of creates a predicate that can be used to compare a value with predefined expected values
func Of[T comparable](expected ...T) predicate.Predicate[T] {
	return as.ErrTail(set.Of(expected...).Contains)
}
