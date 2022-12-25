// common operations of c.Iterable based collections
package op

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/ptr"
)

// Map instantiates Iterator that converts elements with a converter and returns them
func Map[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, by c.Converter[From, To]) c.Iterator[To] {
	return it.Map(elements.Begin(), by)
}

// MapFit additionally filters 'From' elements
func MapFit[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	return it.MapFit(elements.Begin(), fit, by)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, by c.Flatter[From, To]) c.Iterator[To] {
	return ptr.Of(it.Flatt(elements.Begin(), by))
}

// FlattFit additionally filters 'From' elements
func FlattFit[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	return ptr.Of(it.FlattFit(elements.Begin(), fit, flatt))
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT c.Iterable[c.Iterator[T]]](elements IT, filter c.Predicate[T]) c.Iterator[T] {
	return it.Filter(elements.Begin(), filter)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterable[c.Iterator[*T]]](elements IT) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

// Reduce reduces elements to an one
func Reduce[T any, IT c.Iterable[c.Iterator[T]]](elements IT, by c.Binary[T]) T {
	return it.Reduce(elements.Begin(), by)
}

// ToSlice converts an Iterable to a slice
func ToSlice[T any, IT c.Iterable[c.Iterator[T]]](elements IT) []T {
	return it.ToSlice[T](elements.Begin())
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[IT], IT c.Iterator[T]](elements C, by c.Converter[T, K]) c.MapPipe[K, T, map[K][]T] {
	return it.Group(elements.Begin(), by)
}
