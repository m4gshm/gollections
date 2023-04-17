// Package it provides implementations of untility methods for the c.Iterator
package it

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/kvit/group"
	"github.com/m4gshm/gollections/notsafe"
)

// Convert instantiates Iterator that converts elements with a converter and returns them.
func Convert[From, To any, IT any](elements IT, next func() (From, bool), by func(From) To) ConvertIter[From, To, IT] {
	return ConvertIter[From, To, IT]{iter: elements, next: next, by: by}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To, IT any](iter IT, next func() (From, bool), filter func(From) bool, by func(From) To) ConvertFitIter[From, To, IT] {
	return ConvertFitIter[From, To, IT]{iter: iter, next: next, by: by, filter: filter}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To, IT any](iter IT, next func() (From, bool), by func(From) []To) Flatten[From, To, IT] {
	return Flatten[From, To, IT]{iter: iter, next: next, flatt: by, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FilterAndFlatt additionally filters 'From' elements.
func FilterAndFlatt[From, To, IT any](iter IT, next func() (From, bool), filter func(From) bool, flatt func(From) []To) FlattenFit[From, To, IT] {
	return FlattenFit[From, To, IT]{iter: iter, next: next, flatt: flatt, filter: filter, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Filter creates an Iterator that checks elements by filters and returns successful ones.
func Filter[T, IT any](iter IT, next func() (T, bool), filter func(T) bool) Fit[T, IT] {
	return Fit[T, IT]{iter: iter, next: next, by: filter}
}

// NotNil creates an Iterator that filters nullable elements.
func NotNil[T any, IT c.Iterator[*T]](elements IT) Fit[*T, IT] {
	return Filter(elements, elements.Next, check.NotNil[T])
}

// ConvertKV creates an Iterator that applies a transformer to iterable key\values.
func ConvertKV[K, V any, IT c.KVIterator[K, V], k2, v2 any](elements IT, by func(K, V) (k2, v2)) ConvertKVIter[K, V, IT, k2, v2, func(K, V) (k2, v2)] {
	return ConvertKVIter[K, V, IT, k2, v2, func(K, V) (k2, v2)]{iter: elements, by: by}
}

// FilterKV creates an Iterator that checks elements by a filter and returns successful ones
func FilterKV[K, V any, IT c.KVIterator[K, V]](elements IT, filter func(K, V) bool) FitKV[K, V, IT] {
	return FitKV[K, V, IT]{iter: elements, by: filter}
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, IT c.Iterator[T]](elements IT, keyExtractor func(T) K) c.MapPipe[K, T, map[K][]T] {
	return GroupAndConvert(elements, keyExtractor, as.Is[T])
}

// GroupAndConvert transforms iterable elements to the MapPipe based on applying key extractor to the elements
func GroupAndConvert[T any, K comparable, V any, IT c.Iterator[T]](elements IT, keyExtractor func(T) K, valueConverter func(T) V) c.MapPipe[K, V, map[K][]V] {
	return NewKVPipe(NewKeyValuer(elements, elements.Next, keyExtractor, valueConverter), group.Of[K, V])
}

// ToSlice converts an Iterator to a slice
func ToSlice[T any, IT c.Iterator[T]](elements IT) []T {
	var s []T
	a := any(elements)

	if sized, ok := a.(c.Sized); !ok {
		s = make([]T, 0)
	} else if cap := sized.Cap(); cap > 0 {
		s = make([]T, 0, cap)
	} else {
		s = make([]T, 0)
	}
	for v, ok := elements.Next(); ok; v, ok = elements.Next() {
		s = append(s, v)
	}
	return s
}

// Reduce reduces elements to an one
func Reduce[T any](next func() (T, bool), by func(T, T) T) T {
	var result T
	if v, ok := next(); ok {
		result = v
	} else {
		return result
	}
	for v, ok := next(); ok; v, ok = next() {
		result = by(result, v)
	}
	return result
}

// ReduceKV reduces key/values elements to an one
func ReduceKV[K, V any](next func() (K, V, bool), by func(K, V, K, V) (K, V)) (K, V) {
	var rk K
	var rv V
	if k, v, ok := next(); ok {
		rk, rv = k, v
	} else {
		return rk, rv
	}
	for k, v, ok := next(); ok; k, v, ok = next() {
		rk, rv = by(rk, rv, k, v)
	}
	return rk, rv
}
