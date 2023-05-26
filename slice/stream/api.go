// Package stream provides helper functions for transforming a slice to a stream
package stream

import (
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/slice/iter"
	"github.com/m4gshm/gollections/stream"
)

// Convert returns a stream that applies the 'converter' function to the 'elements' slice
func Convert[FS ~[]From, From, To any](elements FS, by func(From) To) stream.Iter[To] {
	conv := iter.Convert(elements, by)
	return stream.New(conv.Next)
}

// FilterAndConvert returns a stream that filters source elements and converts them
func FilterAndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	f := iter.FilterAndConvert(elements, filter, converter)
	return stream.New(f.Next)
}

// Flat returns a stream that converts the collection elements into slices and then flattens them to one level
func Flat[FS ~[]From, From, To any](elements FS, by func(From) []To) stream.Iter[To] {
	f := iter.Flat(elements, by)
	return stream.New(f.Next)
}

// FilterAndFlat filters source elements and extracts slices of 'To' by the 'flattener' function
func FilterAndFlat[FS ~[]From, From, To any](elements FS, filter func(From) bool, flattener func(From) []To) stream.Iter[To] {
	f := iter.FilterAndFlat(elements, filter, flattener)
	return stream.New(f.Next)
}

// Filter instantiates a stream that checks elements by the 'filter' function and returns successful ones
func Filter[TS ~[]T, T any](elements TS, filter func(T) bool) stream.Iter[T] {
	f := iter.Filter(elements, filter)
	return stream.New(f.Next)
}

// NotNil instantiates a stream that filters nullable elements
func NotNil[T any, TRS ~[]*T](elements TRS) stream.Iter[*T] {
	return Filter(elements, not.Nil[T])
}
