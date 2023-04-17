// Package iterable consists of common operations of c.Iterable based collections
package iterable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, by func(From) To) c.Pipe[To, []To] {
	b := elements.Begin()
	return iter.NewPipe[To](iter.Convert(b, b.Next, by))
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, filter func(From) bool, converter func(From) To) c.Pipe[To, []To] {
	b := elements.Begin()
	return iter.NewPipe[To](iter.FilterAndConvert(b, b.Next, filter, converter))
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flatter from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, by func(From) []To) c.Pipe[To, []To] {
	b := elements.Begin()
	f := iter.Flatt(b, b.Next, by)
	return iter.NewPipe[To](&f)
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, filter func(From) bool, flatt func(From) []To) c.Pipe[To, []To] {
	b := elements.Begin()
	f := iter.FilterAndFlatt(b, b.Next, filter, flatt)
	return iter.NewPipe[To](&f)
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT c.Iterable[c.Iterator[T]]](elements IT, filter func(T) bool) c.Pipe[T, []T] {
	b := elements.Begin()
	return iter.NewPipe[T](iter.Filter(b, b.Next, filter))
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterable[c.Iterator[*T]]](elements IT) c.Pipe[*T, []*T] {
	return Filter(elements, check.NotNil[T])
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[IT], IT c.Iterator[T]](elements C, by func(T) K) c.MapPipe[K, T, map[K][]T] {
	return iter.Group(elements.Begin(), by)
}
