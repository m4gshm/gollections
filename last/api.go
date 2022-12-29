package last

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/last"
)

// Of an alias of the slice.Last
func Of[T any](elements ...T) OfElements[T] {
	return OfElements[T]{elements: elements}
}

// By an alias of the slice.Last
func By[T any](by c.Predicate[T]) ByPredicate[T] {
	return ByPredicate[T]{by: by}
}

// ByPredicate tail of the By method
type ByPredicate[T any] struct {
	by c.Predicate[T]
}

func (l ByPredicate[T]) Of(elements ...T) (T, bool) {
	return last.Of(elements, l.by)
}

// OfElements tail of the Of method
type OfElements[T any] struct {
	elements []T
}

func (l OfElements[T]) By(by c.Predicate[T]) (T, bool) {
	return last.Of(l.elements, by)
}
