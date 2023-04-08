package more

import (
	"github.com/m4gshm/gollections/predicate"
	"golang.org/x/exp/constraints"
)

func Than[T constraints.Ordered](min T) predicate.Predicate[T] {
	return func(v T) bool { return v > min }
}

func OrEq[T constraints.Ordered](min T) predicate.Predicate[T] {
	return func(v T) bool { return v >= min }
}
