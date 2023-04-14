package slice

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/slice"
	"github.com/m4gshm/gollections/ptr"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[FS ~[]From, From, To any](elements FS, by func(From) To) c.Iterator[To] {
	return ptr.Of(slice.Convert(elements, by))
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, by func(From) To) c.Iterator[To] {
	return ptr.Of(slice.FilterAndConvert(elements, filter, by))
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[FS ~[]From, From, To any](elements FS, by func(From) []To) c.Iterator[To] {
	return ptr.Of(slice.Flatt(elements, by))
}

// FilterAndFlatt additionally filters 'From' elements.
func FilterAndFlatt[FS ~[]From, From, To any](elements FS, filter func(From) bool, flatt func(From) []To) c.Iterator[To] {
	return ptr.Of(slice.FilterAndFlatt(elements, filter, flatt))
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones.
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) c.Iterator[T] {
	return ptr.Of(slice.Filter(elements, filter))
}

// NotNil instantiates Iterator that filters nullable elements.
func NotNil[T any, TRS ~[]*T](elements TRS) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}
