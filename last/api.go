package last

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/last"
)

// Of an alias of the slice.Last
func Of[T any](elements ...T) ofElements[T] {
	return ofElements[T]{elements: elements}
}

// By an alias of the slice.Last
func By[T any](by c.Predicate[T]) byPredicate[T] {
	return byPredicate[T]{by: by}
}

type byPredicate[T any] struct {
	by c.Predicate[T]
}

func (l byPredicate[T]) Of(elements ...T) (T, bool) {
	return last.Of(elements, l.by)
}

type ofElements[T any] struct {
	elements []T
}

func (l ofElements[T]) By(by c.Predicate[T]) (T, bool) {
	return last.Of(l.elements, by)
}
