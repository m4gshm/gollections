package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/slice"
)

// Map instantiates Iterator that converts elements with a converter and returns them
func Map[From, To any, FS ~[]From](elements FS, by c.Converter[From, To]) c.Iterator[To] {
	iter := slice.Map(elements, by)
	return &iter
}

// MapFit additionally filters 'From' elements.
func MapFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], by c.Converter[From, To]) c.Iterator[To] {
	iter := slice.MapFit(elements, fit, by)
	return &iter
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any, FS ~[]From](elements FS, by c.Flatter[From, To]) c.Iterator[To] {
	iter := slice.Flatt(elements, by)
	return &iter
}

// FlattFit additionally filters 'From' elements.
func FlattFit[From, To any, FS ~[]From](elements FS, fit c.Predicate[From], flatt c.Flatter[From, To]) c.Iterator[To] {
	iter := slice.FlattFit(elements, fit, flatt)
	return &iter
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[T any, TS ~[]T](elements TS, filter c.Predicate[T]) c.Iterator[T] {
	iter := slice.Filter(elements, filter)
	return &iter
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}
