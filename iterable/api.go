// Package iterable consists of common operations of c.Iterable based collections
package iterable

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/check"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
)

// Convert instantiates Iterator that converts elements with a converter and returns them
func Convert[I c.Iterator[From], From, To any, Col c.Iterable[I]](elements Col, converter func(From) To) loop.StreamIter[To] {
	b := elements.Begin()
	return loop.Stream(loop.Convert(b.Next, converter).Next)
}

// FilterAndConvert additionally filters 'From' elements
func FilterAndConvert[From, To any, Col c.Iterable[c.Iterator[From]]](elements Col, filter func(From) bool, converter func(From) To) loop.StreamIter[To] {
	b := elements.Begin()
	f := loop.FilterAndConvert(b.Next, filter, converter)
	return loop.Stream(f.Next)
}

// Flatt instantiates Iterator that extracts slices of 'To' by a flattener from elements of 'From' and flattens as one iterable collection of 'To' elements
func Flatt[From, To any, Col c.Iterable[c.Iterator[From]]](elements Col, by func(From) []To) loop.StreamIter[To] {
	b := elements.Begin()
	f := loop.Flatt(b.Next, by)
	return loop.Stream(f.Next)
}

// FilterAndFlatt additionally filters 'From' elements
func FilterAndFlatt[From, To any, Col c.Iterable[c.Iterator[From]]](elements Col, filter func(From) bool, flatt func(From) []To) loop.StreamIter[To] {
	b := elements.Begin()
	f := loop.FilterAndFlatt(b.Next, filter, flatt)
	return loop.Stream(f.Next)
}

// Filter instantiates Iterator that checks elements by filters and returns successful ones
func Filter[I c.Iterator[T], T any, Col c.Iterable[I]](elements Col, filter func(T) bool) loop.StreamIter[T] {
	b := elements.Begin()
	f := loop.Filter(b.Next, filter)
	return loop.Stream(f.Next)
}

// NotNil instantiates Iterator that filters nullable elements
func NotNil[T any, Col c.Iterable[c.Iterator[*T]]](elements Col) loop.StreamIter[*T] {
	return Filter[c.Iterator[*T]](elements, check.NotNil[T])
}

// Group groups elements to slices by a converter and returns a map
func Group[T any, K comparable, C c.Iterable[c.Iterator[T]]](elements C, by func(T) K) kvloop.StreamIter[K, T, map[K][]T] {
	it := loop.NewKeyValuer(elements.Begin().Next, by, as.Is[T])
	return kvloop.Stream(it.Next, kvloop.Group[K, T])
}
