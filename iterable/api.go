// Package iterable consists of common operations of c.Iterable based collections
package iterable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/iterable/transform"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterable[From]](elements IT, converter func(From) To) c.Transform[To] {
	b := elements.Begin()
	return transform.New(iter.Convert(b.Next, converter).Next)
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, converter func(From) To) c.Transform[To] {
	b := elements.Begin()
	f := iter.FilterAndConvert(b.Next, filter, converter)
	return transform.New(f.Next)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flattener from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterable[From]](elements IT, by func(From) []To) c.Transform[To] {
	b := elements.Begin()
	f := iter.Flatt(b.Next, by)
	return transform.New(f.Next)
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, flatt func(From) []To) c.Transform[To] {
	b := elements.Begin()
	f := iter.FilterAndFlatt(b.Next, filter, flatt)
	return transform.New(f.Next)
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT c.Iterable[T]](elements IT, filter func(T) bool) c.Transform[T] {
	b := elements.Begin()
	f := iter.Filter(b.Next, filter)
	return transform.New(f.Next)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterable[*T]](elements IT) c.Transform[*T] {
	return Filter(elements, check.NotNil[T])
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[T]](elements C, by func(T) K) c.MapTransform[K, T, map[K][]T] {
	return iter.Group(elements.Begin().Next, by)
}
