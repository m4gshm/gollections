// Package iter provides implementations of untility methods for the c.Iterator
package iter

import (
	"github.com/m4gshm/gollections/check"
	"github.com/m4gshm/gollections/notsafe"
)

// Convert instantiates Iterator that converts elements with a converter and returns them.
func Convert[From, To any](next func() (From, bool), converter func(From) To) ConvertIter[From, To] {
	return ConvertIter[From, To]{next: next, converter: converter}
}

// FilterAndConvert additionally filters 'From' elements.
func FilterAndConvert[From, To any](next func() (From, bool), filter func(From) bool, by func(From) To) ConvertFitIter[From, To] {
	return ConvertFitIter[From, To]{next: next, by: by, filter: filter}
}

// Flatt instantiates Iterator that extracts slices of 'To' by a Flattener from elements of 'From' and flattens as one iterable collection of 'To' elements.
func Flatt[From, To any](next func() (From, bool), converter func(From) []To) Flatten[From, To] {
	return Flatten[From, To]{next: next, flatt: converter, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// FilterAndFlatt additionally filters 'From' elements.
func FilterAndFlatt[From, To any](next func() (From, bool), filter func(From) bool, flatt func(From) []To) FlattenFit[From, To] {
	return FlattenFit[From, To]{next: next, flatt: flatt, filter: filter, elemSizeTo: notsafe.GetTypeSize[To]()}
}

// Filter creates an Iterator that checks elements by filters and returns successful ones.
func Filter[T any](next func() (T, bool), filter func(T) bool) Fit[T] {
	return Fit[T]{next: next, by: filter}
}

// NotNil creates an Iterator that filters nullable elements.
func NotNil[T any](next func() (*T, bool)) Fit[*T] {
	return Filter(next, check.NotNil[T])
}
