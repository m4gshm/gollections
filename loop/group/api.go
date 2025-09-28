// Package group provides short aliases for functions that are used to group elements retrieved by a seq
package group

import (
	"github.com/m4gshm/gollections/loop"
)

// Of is a short alias for loop.Group
func Of[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	return loop.Group(next, keyExtractor, valExtractor)
}

// ByMultiple is a short alias for loop.GroupByMultiple
func ByMultiple[T any, K comparable, V any](next func() (T, bool), keysExtractor func(T) []K, valsExtractor func(T) []V) map[K][]V {
	return loop.GroupByMultiple(next, keysExtractor, valsExtractor)
}

// ByMultipleKeys is a short alias for loop.GroupByMultipleKeys
func ByMultipleKeys[T any, K comparable, V any](next func() (T, bool), keysExtractor func(T) []K, valExtractor func(T) V) map[K][]V {
	return loop.GroupByMultipleKeys(next, keysExtractor, valExtractor)
}

// ByMultipleValues is a short alias for loop.GroupByMultipleVals
func ByMultipleValues[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valsExtractor func(T) []V) map[K][]V {
	return loop.GroupByMultipleValues(next, keyExtractor, valsExtractor)
}
