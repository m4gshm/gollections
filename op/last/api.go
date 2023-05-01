// Package last provides helpers for retrieving a last element of a slice that satisfies a condition
package last

import (
	"github.com/m4gshm/gollections/slice/last"
)

// Of an alias of the slice.Last
func Of[T any](elements ...T) OfElements[T] {
	return OfElements[T]{elements: elements}
}

// By an alias of the slice.Last
func By[T any](by func(T) bool) ByPredicate[T] {
	return ByPredicate[T]{predicate: by}
}

// ByPredicate tail of the By method
type ByPredicate[T any] struct {
	predicate func(T) bool
}

// Of the predicate apply method
func (l ByPredicate[T]) Of(elements ...T) (T, bool) {
	return last.Of(elements, l.predicate)
}

// OfElements tail of the Of method
type OfElements[T any] struct {
	elements []T
}

// By the predicate apply method
func (l OfElements[T]) By(predicate func(T) bool) (T, bool) {
	return last.Of(l.elements, predicate)
}
