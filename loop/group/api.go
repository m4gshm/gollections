package group

import (
	"github.com/m4gshm/gollections/loop"
)

// Of is a short alias for loop.Group
func Of[T any, K comparable](next func() (T, bool), keyProducer func(T) K) map[K][]T {
	return loop.Group(next, keyProducer)
}

// InMultiple is a short alias for loop.GroupInMultiple
func InMultiple[T any, K comparable](next func() (T, bool), keysProducer func(T) []K) map[K][]T {
	return loop.GroupInMultiple(next, keysProducer)
}
