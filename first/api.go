// Package first provides helpers for retrieving a first element of a slice that satisfies a condition
package first

import (
	"github.com/m4gshm/gollections/slice/first"
)

// Of the first part of an expression first.Of(elements...).By(tester)
func Of[T any](elements ...T) OfElements[T] {
	return OfElements[T]{elements: elements}
}

// By the first part of an expression first.By(tester).Of(elements...)
func By[T any](by func(T) bool) ByPredicate[T] {
	return ByPredicate[T]{by: by}
}

// ByPredicate is tail prducer of the first.By
type ByPredicate[T any] struct {
	by func(T) bool
}

// Of the finish part of an expression first.By(tester).Of(elements...)
func (l ByPredicate[T]) Of(elements ...T) (T, bool) {
	return first.Of(elements, l.by)
}

// OfElements is tail prducer of the first.Of
type OfElements[T any] struct {
	elements []T
}

// By the finish part of an expression first.Of(elements...).By(tester)
func (l OfElements[T]) By(by func(T) bool) (T, bool) {
	return l.Where(by)
}

// Where the finish part of an expression first.Of(elements...).Where(condition)
func (l OfElements[T]) Where(condition func(T) bool) (T, bool) {
	return first.Of(l.elements, condition)
}
