package slice

import (
	"github.com/m4gshm/container/check"
	"github.com/m4gshm/container/conv"
	impl "github.com/m4gshm/container/slice/impl/slice"
	"github.com/m4gshm/container/typ"
)

func Of[T any](elements ...T) []T {
	return elements
}

//Map creates a lazy Iterator that converts elements with a converter and returns them
func Map[From, To any](elements []From, by conv.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

// additionally filters 'From' elements by filters
func MapFit[From, To any](elements []From, fit check.Predicate[From], by conv.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any](elements []From, by conv.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

// additionally checks 'From' elements by fit
func FlattFit[From, To any](elements []From, fit check.Predicate[From], flatt conv.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones
func Filter[T any](elements []T, filter check.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates a lazy Iterator that filters nullable elements
func NotNil[T any](elements []*T) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}
