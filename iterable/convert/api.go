// Package convert provides converation helpers for collection implementations
package convert

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/not"
)

// AndConvert - convert.AndConvert makes double converts From->Intermediate->To of the elements
func AndConvert[From, To, Too any](elements c.Iterable[From], firsConverter func(From) To, secondConverter func(To) Too) c.Stream[Too] {
	return iterable.Convert(iterable.Convert(elements, firsConverter), secondConverter)
}

// AndFilter - convert.AndFilter converts only filtered elements and returns them
func AndFilter[From, To any](elements c.Iterable[From], converter func(From) To, filter func(To) bool) c.Stream[To] {
	return iterable.Filter(iterable.Convert(elements, converter), filter)
}

// NotNil - convert.NotNil converts only not nil elements and returns them
func NotNil[From, To any](elements c.Iterable[*From], converter func(*From) To) c.Stream[To] {
	return iterable.FilterAndConvert(elements, not.Nil[From], converter)
}
