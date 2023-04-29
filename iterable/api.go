// Package iterable consists of common operations of c.Iterable based collections
package iterable

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	kvstream "github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/stream"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[I c.Iterator[From], From, To any, Col c.Iterable[I]](elements Col, converter func(From) To) stream.Iter[To] {
	b := elements.Begin()
	return stream.New(loop.Convert(b.Next, converter).Next)
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[I c.Iterator[From], From, To any, Col c.Iterable[I]](elements Col, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	b := elements.Begin()
	f := loop.FilterAndConvert(b.Next, filter, converter)
	return stream.New(f.Next)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[I c.Iterator[From], From, To any, Col c.Iterable[I]](elements Col, by func(From) []To) stream.Iter[To] {
	b := elements.Begin()
	f := loop.Flatt(b.Next, by)
	return stream.New(f.Next)
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[I c.Iterator[From], From, To any, Col c.Iterable[I]](elements Col, filter func(From) bool, flatt func(From) []To) stream.Iter[To] {
	b := elements.Begin()
	f := loop.FilterAndFlatt(b.Next, filter, flatt)
	return stream.New(f.Next)
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[I c.Iterator[T], T any, Col c.Iterable[I]](elements Col, filter func(T) bool) stream.Iter[T] {
	b := elements.Begin()
	f := loop.Filter(b.Next, filter)
	return stream.New(f.Next)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[I c.Iterator[*T], T any, Col c.Iterable[I]](elements Col) stream.Iter[*T] {
	return Filter[I](elements, check.NotNil[T])
}

// Group groups elements to slices by a converter and returns a map
func Group[I c.Iterator[T], T any, K comparable, C c.Iterable[I]](elements C, by func(T) K) kvstream.Iter[K, T, map[K][]T] {
	it := loop.NewKeyValuer(elements.Begin().Next, by, as.Is[T])
	return kvstream.New(it.Next, kvloop.Group[K, T])
}
