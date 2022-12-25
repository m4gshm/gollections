package it

import (
	"errors"

	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/kvit/group"
	"github.com/m4gshm/gollections/notsafe"
)

// ErrBreak is For, Track breaker
var ErrBreak = errors.New("Break")

// Map instantiates Iterator that converts elements with a converter and returns them
func Map[From, To any, IT c.Iterator[From]](elements IT, by c.Converter[From, To]) ConvertIter[From, To, IT, c.Converter[From, To]] {
	return ConvertIter[From, To, IT, c.Converter[From, To]]{iter: elements, by: by}
}

// MapFit additionally filters 'From' elements.
func MapFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], by c.Converter[From, To]) ConvertFitIter[From, To, IT] {
	return ConvertFitIter[From, To, IT]{iter: elements, by: by, fit: fit}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, IT c.Iterator[From]](elements IT, by c.Flatter[From, To]) Flatten[From, To, IT] {
	return Flatten[From, To, IT]{iter: elements, flatt: by, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, IT c.Iterator[From]](elements IT, fit c.Predicate[From], flatt c.Flatter[From, To]) FlattenFit[From, To, IT] {
	return FlattenFit[From, To, IT]{iter: elements, flatt: flatt, fit: fit, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[T any, IT c.Iterator[T]](elements IT, filter c.Predicate[T]) Fit[T, IT] {
	return Fit[T, IT]{iter: elements, by: filter}
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, IT c.Iterator[*T]](elements IT) Fit[*T, IT] {
	return Filter(elements, check.NotNil[T])
}

// MapKV instantiates Iterator that converts elements with a converter and returns them
func MapKV[K, V any, IT c.KVIterator[K, V], k2, v2 any](elements IT, by c.BiConverter[K, V, k2, v2]) ConvertKVIter[K, V, IT, k2, v2, c.BiConverter[K, V, k2, v2]] {
	return ConvertKVIter[K, V, IT, k2, v2, c.BiConverter[K, V, k2, v2]]{iter: elements, by: by}
}

// FilterKV instantiates Iterator that checks elements by filters and returns successful ones.
func FilterKV[K, V any, IT c.KVIterator[K, V]](elements IT, filter c.BiPredicate[K, V]) FitKV[K, V, IT] {
	return FitKV[K, V, IT]{iter: elements, by: filter}
}

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable, IT c.Iterator[T]](elements IT, by c.Converter[T, K]) c.MapPipe[K, T, map[K][]T] {
	return GroupMap(elements, by, as.Is[T])
}

// GroupMap transforms iterable elements to the MapPipe based on applying key extractor to the elements
func GroupMap[T any, K comparable, V any, IT c.Iterator[T]](elements IT, by c.Converter[T, K], val c.Converter[T, V]) c.MapPipe[K, V, map[K][]V] {
	return NewKVPipe(NewKeyValuer(elements, by, val), group.Of[K, V])
}

// For applies a walker to elements of an Iterator. To stop walking just return the ErrBreak.
func For[T any, IT c.Iterator[T]](elements IT, walker func(T) error) error {
	for v, ok := elements.Next(); ok; v, ok = elements.Next() {
		if err := walker(v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies a walker to elements of an Iterator.
func ForEach[T any, IT c.Iterator[T]](elements IT, walker func(T)) {
	for v, ok := elements.Next(); ok; v, ok = elements.Next() {
		walker(v)
	}
}

// ForEachFit applies a walker to elements that satisfy a predicate condition.
func ForEachFit[T any, IT c.Iterator[T]](elements IT, walker func(T), fit c.Predicate[T]) {
	for v, ok := elements.Next(); ok && fit(v); v, ok = elements.Next() {
		walker(v)
	}
}

// Reduce reduces elements to an one.
func Reduce[T any, IT c.Iterator[T]](elements IT, by c.Binary[T]) T {
	var result T
	if v, ok := elements.Next(); ok {
		result = v
	} else {
		return result
	}
	for v, ok := elements.Next(); ok; v, ok = elements.Next() {
		result = by(result, v)
	}
	return result
}

// ReduceKV reduces key/values elements to an one.
func ReduceKV[K, V any, IT c.KVIterator[K, V]](elements IT, by c.Quaternary[K, V]) (K, V) {
	var rk K
	var rv V
	if k, v, ok := elements.Next(); ok {
		rk, rv = k, v
	} else {
		return rk, rv
	}
	for k, v, ok := elements.Next(); ok; k, v, ok = elements.Next() {
		rk, rv = by(rk, rv, k, v)
	}
	return rk, rv
}

// ToSlice converts an Iterator to a slice.
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

// First returns the first element that satisfies requirements of the predicate 'fit'
func First[T any, IT c.Iterator[T]](elements IT, fit c.Predicate[T]) (T, bool) {
	for one, ok := elements.Next(); ok; one, ok = elements.Next() {
		if fit(one) {
			return one, true
		}
	}
	var no T
	return no, false
}
