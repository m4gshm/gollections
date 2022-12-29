package first

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/first"
)

// Of the first part of an expression first.Of(elements...).By(tester)
func Of[T any](elements ...T) OfElements[T] {
	return OfElements[T]{elements: elements}
}

// By the first part of an expression first.By(tester).Of(elements...)
func By[T any](by c.Predicate[T]) ByPredicate[T] {
	return ByPredicate[T]{by: by}
}

// ByPredicate is tail prducer of the first.By
type ByPredicate[T any] struct {
	by c.Predicate[T]
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
func (l OfElements[T]) By(by c.Predicate[T]) (T, bool) {
	return first.Of(l.elements, by)
}
