// Package op consists of common operations of c.Iterable based collections
package op

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/ptr"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, by func(From) To) c.Iterator[To] {
	b := elements.Begin()
	return it.Convert(b, b.Next, by)
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, filter func(From) bool, converter func(From) To) c.Iterator[To] {
	b := elements.Begin()
	return it.FilterAndConvert(b, b.Next, filter, converter)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, by func(From) []To) c.Iterator[To] {
	b := elements.Begin()
	return ptr.Of(it.Flatt(b, b.Next, by))
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, filter func(From) bool, flatt func(From) []To) c.Iterator[To] {
	b := elements.Begin()
	return ptr.Of(it.FilterAndFlatt(b, b.Next, filter, flatt))
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT c.Iterable[c.Iterator[T]]](elements IT, filter func(T) bool) c.Iterator[T] {
	b := elements.Begin()
	return it.Filter(b, b.Next, filter)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterable[c.Iterator[*T]]](elements IT) c.Iterator[*T] {
	return Filter(elements, check.NotNil[T])
}

// Reduce reduces elements to an one
func Reduce[T any, IT c.Iterable[c.Iterator[T]]](elements IT, by func(T, T) T) T {
	return it.Reduce(elements.Begin().Next, by)
}

// ToSlice converts an Iterable to a slice
func ToSlice[T any, IT c.Iterable[c.Iterator[T]]](elements IT) []T {
	return it.ToSlice[T](elements.Begin())
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[IT], IT c.Iterator[T]](elements C, by func(T) K) c.MapPipe[K, T, map[K][]T] {
	return it.Group(elements.Begin(), by)
}
