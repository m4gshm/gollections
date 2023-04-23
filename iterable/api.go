// Package iterable consists of common operations of c.Iterable based collections
package iterable

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/loop/iter"
	kvstream "github.com/m4gshm/gollections/loop/kv/stream"
	"github.com/m4gshm/gollections/loop/stream"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[From, To any, IT c.Iterable[From]](elements IT, converter func(From) To) c.Stream[To] {
	b := elements.Begin()
	return stream.New(iter.Convert(b.Next, converter).Next)
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, converter func(From) To) c.Stream[To] {
	b := elements.Begin()
	f := iter.FilterAndConvert(b.Next, filter, converter)
	return stream.New(f.Next)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flattener from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, IT c.Iterable[From]](elements IT, by func(From) []To) c.Stream[To] {
	b := elements.Begin()
	f := iter.Flatt(b.Next, by)
	return stream.New(f.Next)
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, flatt func(From) []To) c.Stream[To] {
	b := elements.Begin()
	f := iter.FilterAndFlatt(b.Next, filter, flatt)
	return stream.New(f.Next)
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[T any, IT c.Iterable[T]](elements IT, filter func(T) bool) c.Stream[T] {
	b := elements.Begin()
	f := iter.Filter(b.Next, filter)
	return stream.New(f.Next)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, IT c.Iterable[*T]](elements IT) c.Stream[*T] {
	return Filter(elements, check.NotNil[T])
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[T]](elements C, by func(T) K) c.MapStream[K, T, map[K][]T] {
	return kvstream.Group(elements.Begin().Next, by)
}
