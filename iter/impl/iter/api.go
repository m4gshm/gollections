// Package iter provides implementations of untility methods for the c.Iterator
package iter

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/kviter/group"
	"github.com/m4gshm/gollections/notsafe"
)

// Convert instantiates Iterator that converts elements with a converter and returns them.
func Convert[From, To any, IT any](elements IT, next func() (From, bool), converter func(From) To) ConvertIter[From, To] {
	return ConvertIter[From, To]{next: next, converter: converter}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To any, IT any](iterator IT, next func() (From, bool), filter func(From) bool, by func(From) To) ConvertFitIter[From, To] {
	return ConvertFitIter[From, To]{next: next, by: by, filter: filter}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](next func() (From, bool), converter func(From) []To) Flatten[From, To] {
	return Flatten[From, To]{next: next, flatt: converter, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FilterAndFlatt additionally filters 'From' elements.
func FilterAndFlatt[From, To any](next func() (From, bool), filter func(From) bool, flatt func(From) []To) FlattenFit[From, To] {
	return FlattenFit[From, To]{next: next, flatt: flatt, filter: filter, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Filter creates an Iterator that checks elements by filters and returns successful ones.
func Filter[T any](next func() (T, bool), filter func(T) bool) Fit[T] {
	return Fit[T]{next: next, by: filter}
}

// NotNil creates an Iterator that filters nullable elements.
func NotNil[T any](next func() (*T, bool)) Fit[*T] {
	return Filter(next, check.NotNil[T])
}

// ConvertKV creates an Iterator that applies a transformer to iterable key\values.
func ConvertKV[K, V any, IT c.KVIterator[K, V], k2, v2 any](elements IT, by func(K, V) (k2, v2)) ConvertKVIter[K, V, k2, v2, func(K, V) (k2, v2)] {
	return ConvertKVIter[K, V, k2, v2, func(K, V) (k2, v2)]{next: elements.Next, by: by}
}

// FilterKV creates an Iterator that checks elements by a filter and returns successful ones
func FilterKV[K, V any, IT c.KVIterator[K, V]](elements IT, filter func(K, V) bool) FitKV[K, V] {
	return FitKV[K, V]{next: elements.Next, by: filter}
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](next func() (T, bool), keyExtractor func(T) K) KVIterPipe[K, T, map[K][]T] {
	return GroupAndConvert(next, keyExtractor, as.Is[T])
}

// GroupAndConvert transforms iterable elements to the MapPipe based on applying key extractor to the elements
func GroupAndConvert[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valueConverter func(T) V) KVIterPipe[K, V, map[K][]V] {
	kv := NewKeyValuer(next, keyExtractor, valueConverter)
	return NewKVPipe(&kv, group.Of[K, V])
}
