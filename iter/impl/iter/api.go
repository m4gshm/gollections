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
func Convert[From, To any, IT any](elements IT, next func() (From, bool), by func(From) To) *ConvertIter[From, To, IT] {
	return &ConvertIter[From, To, IT]{iterator: elements, next: next, by: by}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To, IT any](iterator IT, next func() (From, bool), filter func(From) bool, by func(From) To) *ConvertFitIter[From, To, IT] {
	return &ConvertFitIter[From, To, IT]{iterator: iterator, next: next, by: by, filter: filter}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To, IT any](iterator IT, next func() (From, bool), by func(From) []To) Flatten[From, To, IT] {
	return Flatten[From, To, IT]{iterator: iterator, next: next, flatt: by, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FilterAndFlatt additionally filters 'From' elements.
func FilterAndFlatt[From, To, IT any](iterator IT, next func() (From, bool), filter func(From) bool, flatt func(From) []To) FlattenFit[From, To, IT] {
	return FlattenFit[From, To, IT]{iterator: iterator, next: next, flatt: flatt, filter: filter, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Filter creates an Iterator that checks elements by filters and returns successful ones.
func Filter[T, IT any](iterator IT, next func() (T, bool), filter func(T) bool) Fit[T, IT] {
	return Fit[T, IT]{iterator: iterator, next: next, by: filter}
}

// NotNil creates an Iterator that filters nullable elements.
func NotNil[T any, IT c.Iterator[*T]](elements IT) Fit[*T, IT] {
	return Filter(elements, elements.Next, check.NotNil[T])
}

// ConvertKV creates an Iterator that applies a transformer to iterable key\values.
func ConvertKV[K, V any, IT c.KVIterator[K, V], k2, v2 any](elements IT, by func(K, V) (k2, v2)) *ConvertKVIter[K, V, IT, k2, v2, func(K, V) (k2, v2)] {
	return &ConvertKVIter[K, V, IT, k2, v2, func(K, V) (k2, v2)]{iterator: elements, by: by}
}

// FilterKV creates an Iterator that checks elements by a filter and returns successful ones
func FilterKV[K, V any, IT c.KVIterator[K, V]](elements IT, filter func(K, V) bool) FitKV[K, V, IT] {
	return FitKV[K, V, IT]{iterator: elements, by: filter}
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, IT c.Iterator[T]](elements IT, keyExtractor func(T) K) c.MapPipe[K, T, map[K][]T] {
	return GroupAndConvert(elements, keyExtractor, as.Is[T])
}

// GroupAndConvert transforms iterable elements to the MapPipe based on applying key extractor to the elements
func GroupAndConvert[T any, K comparable, V any, IT c.Iterator[T]](elements IT, keyExtractor func(T) K, valueConverter func(T) V) c.MapPipe[K, V, map[K][]V] {
	return NewKVPipe(NewKeyValuer(elements, elements.Next, keyExtractor, valueConverter), group.Of[K, V])
}
