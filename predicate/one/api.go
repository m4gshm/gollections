package one

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/set"
)

func Of[T comparable](expected ...T) c.Predicate[T] {
	return set.Of(expected...).Contains
}
