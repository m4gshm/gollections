// Package iterable consists of common operations of c.Iterable based collections
package iterable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/iter/impl/iter"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterable[From]](elements IT, converter func(From) To) c.Pipe[To] {
	b := elements.Begin()
	return iter.NewPipe(iter.Convert(b.Next, converter).Next)
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, converter func(From) To) c.Pipe[To] {
	b := elements.Begin()
	f := iter.FilterAndConvert(b.Next, filter, converter)
	return iter.NewPipe(f.Next)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flattener from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterable[From]](elements IT, by func(From) []To) c.Pipe[To] {
	b := elements.Begin()
	f := iter.Flatt(b.Next, by)
	return iter.NewPipe(f.Next)
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, flatt func(From) []To) c.Pipe[To] {
	b := elements.Begin()
	f := iter.FilterAndFlatt(b.Next, filter, flatt)
	return iter.NewPipe(f.Next)
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT c.Iterable[T]](elements IT, filter func(T) bool) c.Pipe[T] {
	b := elements.Begin()
	f := iter.Filter(b.Next, filter)
	return iter.NewPipe(f.Next)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterable[*T]](elements IT) c.Pipe[*T] {
	return Filter(elements, check.NotNil[T])
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[T]](elements C, by func(T) K) c.MapPipe[K, T, map[K][]T] {
	return iter.Group(elements.Begin().Next, by)
}
