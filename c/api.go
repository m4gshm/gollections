package c

import (
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it"
	impl "github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

//Map creates a lazy Iterator that converts elements with a converter and returns them.
func Map[From, To any, IT typ.Iterable[typ.Iterator[From]]](elements IT, by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.Map(elements.Begin(), by)
}

// additionally filters 'From' elements by filters.
func MapFit[From, To any, IT typ.Iterable[typ.Iterator[From]]](elements IT, fit typ.Predicate[From], by typ.Converter[From, To]) typ.Iterator[To] {
	return impl.MapFit(elements.Begin(), fit, by)
}

//Flatt creates a lazy Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, IT typ.Iterable[typ.Iterator[From]]](elements IT, by typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.Flatt(elements.Begin(), by)
}

// additionally checks 'From' elements by fit.
func FlattFit[From, To any, IT typ.Iterable[typ.Iterator[From]]](elements IT, fit typ.Predicate[From], flatt typ.Flatter[From, To]) typ.Iterator[To] {
	return impl.FlattFit(elements.Begin(), fit, flatt)
}

//Filter creates a lazy Iterator that checks elements by filters and returns successful ones.
func Filter[T any, IT typ.Iterable[typ.Iterator[T]]](elements IT, filter typ.Predicate[T]) typ.Iterator[T] {
	return impl.Filter(elements.Begin(), filter)
}

//NotNil creates a lazy Iterator that filters nullable elements.
func NotNil[T any, IT typ.Iterable[typ.Iterator[*T]]](elements IT) typ.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

//Reduce reduces elements to an one.
func Reduce[T any, IT typ.Iterable[typ.Iterator[T]]](elements IT, by op.Binary[T]) T {
	return impl.Reduce(elements.Begin(), by)
}

//Slice converts Iterator to slice.
func Slice[T any, IT typ.Iterable[typ.Iterator[T]]](elements IT) []T {
	return impl.Slice[T](elements.Begin())
}

//Group groups elements to slices by a converter and returns map.
func Group[T any, K comparable, IT typ.Iterable[typ.Iterator[T]]](elements IT, by typ.Converter[T, K]) typ.MapPipe[K, T, map[K][]T] {
	return it.Group(elements.Begin(), by)
}
