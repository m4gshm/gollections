package one

import (
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/predicate"
)

func Of[T comparable](expected ...T) predicate.Predicate[T] {
	return set.Of(expected...).Contains
}
