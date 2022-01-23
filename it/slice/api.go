package slice

import (
	"github.com/m4gshm/gollections/check"
	impl "github.com/m4gshm/gollections/it/impl/slice"
	"github.com/m4gshm/gollections/typ"
)

//Map creates the Iterator that converts elements with a converter and returns them.
func Map[From, To any](elements []From, by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements, by)
}

//MapFit additionally filters 'From' elements.
func MapFit[From, To any](elements []From, fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements, fit, by)
}

//Flatt creates the Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](elements []From, by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements, by)
}

//FlattFit additionally filters 'From' elements.
func FlattFit[From, To any](elements []From, fit typ.Predicate[From], flatt typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements, fit, flatt)
}

//Filter creates the Iterator that checks elements by filters and returns successful ones.
func Filter[T any](elements []T, filter typ.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements, filter)
}

//NotNil creates the Iterator that filters nullable elements.
func NotNil[T any](elements []*T) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}
