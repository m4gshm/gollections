package less

import (
	"github.com/m4gshm/gollections/c"
	"golang.org/x/exp/constraints"
)

func Than[T constraints.Ordered](min T) c.Predicate[T] {
	return func(v T) bool { return v < min }
}

func OrEq[T constraints.Ordered](min T) c.Predicate[T] {
	return func(v T) bool { return v <= min }
}
