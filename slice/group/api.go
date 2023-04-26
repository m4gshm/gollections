// Package group provides short aliases for grouping functions
package group

import (
	"github.com/m4gshm/gollections/slice"
)

// Of is a short alias for slice.Group
func Of[TS ~[]T, T any, K comparable, V any](elements TS, keyProducer func(T) K, valProducer func(T) V) map[K][]V {
	return slice.Group(elements, keyProducer, valProducer)
}

// ByMultiple is a short alias for slice.GroupByMultiple
func ByMultiple[TS ~[]T, T any, K comparable, V any](elements TS, keysProducer func(T) []K, valsProducer func(T) []V) map[K][]V {
	return slice.GroupByMultiple(elements, keysProducer, valsProducer)
}

// ByMultipleKeys is a short alias for slice.GroupByMultipleKeys
func ByMultipleKeys[TS ~[]T, T any, K comparable, V any](elements TS, keysProducer func(T) []K, valProducer func(T) V) map[K][]V {
	return slice.GroupByMultipleKeys(elements, keysProducer, valProducer)
}

// ByMultipleValues is a short alias for slice.GroupByMultipleVals
func ByMultipleValues[TS ~[]T, T any, K comparable, V any](elements TS, keyProducer func(T) K, valsProducer func(T) []V) map[K][]V {
	return slice.GroupByMultipleValues(elements, keyProducer, valsProducer)
}
