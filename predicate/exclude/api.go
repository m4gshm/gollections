package exclude

import (
	"github.com/m4gshm/gollections/predicate"
	"github.com/m4gshm/gollections/predicate/one"
)

func All[T comparable](excluded ...T) predicate.Predicate[T] {
	return predicate.Not(one.Of(excluded...))
}
