// Package group provides short aliases for functions used to group elements retrieved by a loop
package group

import (
	"github.com/m4gshm/gollections/loop"
)

// Of is a short alias for loop.Group
func Of[T any, K comparable, V any](next func() (T, bool), keyProducer func(T) K, valProducer func(T) V) map[K][]V {
	return loop.Group(next, keyProducer, valProducer)
}

// ByMultiple is a short alias for loop.GroupByMultiple
func ByMultiple[TS ~[]T, T any, K comparable, V any](next func() (T, bool), keysProducer func(T) []K, valsProducer func(T) []V) map[K][]V {
	return loop.GroupByMultiple(next, keysProducer, valsProducer)
}

// ByMultipleKeys is a short alias for loop.GroupByMultipleKeys
func ByMultipleKeys[TS ~[]T, T any, K comparable, V any](next func() (T, bool), keysProducer func(T) []K, valProducer func(T) V) map[K][]V {
	return loop.GroupByMultipleKeys(next, keysProducer, valProducer)
}

// ByMultipleValues is a short alias for loop.GroupByMultipleVals
func ByMultipleValues[TS ~[]T, T any, K comparable, V any](next func() (T, bool), keyProducer func(T) K, valsProducer func(T) []V) map[K][]V {
	return loop.GroupByMultipleValues(next, keyProducer, valsProducer)
}
