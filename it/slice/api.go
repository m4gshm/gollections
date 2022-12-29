package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/slice"
	"github.com/m4gshm/gollections/ptr"
)

// Map instantiates Iterator that converts elements with a converter and returns them
func Map[FS ~[]From, From, To any](elements FS, by c.Converter[From, To]) c.Iterator[To] {
	return ptr.Of(slice.Map(elements, by))
}

// MapFit additionally filters 'From' elements.
func MapFit[FS ~[]From, From, To any](elements FS, fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	return ptr.Of(slice.MapFit(elements, fit, by))
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[FS ~[]From, From, To any](elements FS, by c.Flatter[From, To]) c.Iterator[To] {
	return ptr.Of(slice.Flatt(elements, by))
}

// FlattFit additionally filters 'From' elements.
func FlattFit[FS ~[]From, From, To any](elements FS, fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	return ptr.Of(slice.FlattFit(elements, fit, flatt))
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[TS ~[]T, T any](elements TS, filter c.Predicate[T]) c.Iterator[T] {
	return ptr.Of(slice.Filter(elements, filter))
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}
