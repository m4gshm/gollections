// Package stream provides helper functions for transforming a slice to a stream
package stream

import (
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/stream"
)

// Convert returns a stream that applies the 'converter' function to the 'elements' slice
func Convert[FS ~[]From, From, To any](elements FS, by func(From) To) stream.Iter[To] {
	return stream.New(loop.Convert(loop.Of(elements...), by))
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	return stream.New(loop.FilterAndConvert(loop.Of(elements...), filter, converter))
}

// Flat returns a stream that converts the collection elements into slices and then flattens them to one level
func Flat[FS ~[]From, From, To any](elements FS, by func(From) []To) stream.Iter[To] {
	return stream.New(loop.Flat(loop.Of(elements...), by))
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[FS ~[]From, From, To any](elements FS, filter func(From) bool, flattener func(From) []To) stream.Iter[To] {
	return stream.New(loop.FilterAndFlat(loop.Of(elements...), filter, flattener))
}

// Filter instantiates a stream that checks elements by the 'filter' function and returns successful ones
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) stream.Iter[T] {
	return stream.New(loop.Filter(loop.Of(elements...), filter))
}

// NotNil instantiates a stream that filters nullable elements
func NotNil[T any, TRS ~[]*T](elements TRS) stream.Iter[*T] {
	return Filter(elements, not.Nil[T])
}
