// Package group provides short aliases for grouping functions
package group

import (
	"github.com/m4gshm/gollections/slice"
)

// Of is a short alias for slice.Group
func Of[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) map[K][]V {
	return slice.Group(elements, keyExtractor, valExtractor)
}

// Order is a short alias for slice.GroupOrder
func Order[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valExtractor func(T) V) ([]K, map[K][]V) {
	return slice.GroupOrder(elements, keyExtractor, valExtractor)
}

// ByMultiple is a short alias for slice.GroupByMultiple
func ByMultiple[TS ~[]T, T any, K comparable, V any](elements TS, keysExtractor func(T) []K, valsExtractor func(T) []V) map[K][]V {
	return slice.GroupByMultiple(elements, keysExtractor, valsExtractor)
}

// ByMultipleKeys is a short alias for slice.GroupByMultipleKeys
func ByMultipleKeys[TS ~[]T, T any, K comparable, V any](elements TS, keysExtractor func(T) []K, valExtractor func(T) V) map[K][]V {
	return slice.GroupByMultipleKeys(elements, keysExtractor, valExtractor)
}

// ByMultipleValues is a short alias for slice.GroupByMultipleVals
func ByMultipleValues[TS ~[]T, T any, K comparable, V any](elements TS, keyExtractor func(T) K, valsExtractor func(T) []V) map[K][]V {
	return slice.GroupByMultipleValues(elements, keyExtractor, valsExtractor)
}
