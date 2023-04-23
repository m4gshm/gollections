package loop

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/notsafe"
)

// ErrBreak is the 'break' statement of the For, Track methods
var ErrBreak = c.ErrBreak

// For applies the 'walker' function for the elements retrieved by the 'next' function. Return the c.ErrBreak to stop
func For[T any](next func() (T, bool), walker func(T) error) error {
	for v, ok := next(); ok; v, ok = next() {
		if err := walker(v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// ForEach applies the 'walker' function to the elements retrieved by the 'next' function
func ForEach[T any](next func() (T, bool), walker func(T)) {
	for v, ok := next(); ok; v, ok = next() {
		walker(v)
	}
}

// ForEachFiltered applies the 'walker' function to the elements retrieved by the 'next' function that satisfy the 'predicate' function condition
func ForEachFiltered[T any](next func() (T, bool), walker func(T), predicate func(T) bool) {
	for v, ok := next(); ok && predicate(v); v, ok = next() {
		walker(v)
	}
}

// First returns the first element that satisfies the condition of the 'predicate' function
func First[T any](next func() (T, bool), predicate func(T) bool) (v T, ok bool) {
	for one, ok := next(); ok; one, ok = next() {
		if predicate(one) {
			return one, true
		}
	}
	return v, ok
}

// Track applies the 'tracker' function to position/element pairs retrieved by the 'next' function. Return the c.ErrBreak to stop tracking..
func Track[I, T any](next func() (I, T, bool), tracker func(I, T) error) error {
	for p, v, ok := next(); ok; p, v, ok = next() {
		if err := tracker(p, v); err == ErrBreak {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

// TrackEach applies the 'tracker' function to position/element pairs retrieved by the 'next' function
func TrackEach[I, T any](next func() (I, T, bool), tracker func(I, T)) {
	for p, v, ok := next(); ok; p, v, ok = next() {
		tracker(p, v)
	}
}

// ToSlice collects the elements retrieved by the 'next' function into a slice
func ToSlice[T any](next func() (T, bool)) []T {
	var s []T
	for v, ok := next(); ok; v, ok = next() {
		s = append(s, v)
	}
	return s
}

// Reduce reduces the elements retrieved by the 'next' function into an one using the 'merge' function
func Reduce[T any](next func() (T, bool), merger func(T, T) T) (result T) {
	if v, ok := next(); ok {
		result = v
	} else {
		return result
	}
	for v, ok := next(); ok; v, ok = next() {
		result = merger(result, v)
	}
	return result
}

// HasAny finds the first element that satisfies the 'predicate' function condition and returns true if successful
func HasAny[T any](next func() (T, bool), predicate func(T) bool) bool {
	_, ok := First(next, predicate)
	return ok
}

// Convert instantiates Iterator that converts elements with a converter and returns them.
func Convert[From, To any](next func() (From, bool), converter func(From) To) ConvertIter[From, To] {
	return ConvertIter[From, To]{next: next, converter: converter}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To any](next func() (From, bool), filter func(From) bool, by func(From) To) ConvertFitIter[From, To] {
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

// ToKV transforms iterable elements to key/value iterator based on applying key extractor to the elements
func ToKV[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valueConverter func(T) V) KeyValuer[T, K, V] {
	kv := NewKeyValuer(next, keyExtractor, valueConverter)
	return kv
}
