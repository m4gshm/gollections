// Package convert provides converation helpers for collection implementations
package convert

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/op/check/not"
	"github.com/m4gshm/gollections/stream"
)

// AndConvert - convert.AndConvert makes double converts From->Intermediate->To of the elements
func AndConvert[From, To, Too any, I c.Iterator[From], IT c.Iterable[From, I]](elements IT, firsConverter func(From) To, secondConverter func(To) Too) stream.Iter[Too] {
	return collection.Convert(collection.Convert(elements, firsConverter), secondConverter)
}

// AndFilter - convert.AndFilter converts only filtered elements and returns them
func AndFilter[From, To any, I c.Iterator[From], IT c.Iterable[From, I]](elements IT, converter func(From) To, filter func(To) bool) stream.Iter[To] {
	return collection.Filter(collection.Convert(elements, converter), filter)
}

// NotNil - convert.NotNil converts only not nil elements and returns them
func NotNil[From, To any, I c.Iterator[*From], IT c.Iterable[*From, I]](elements IT, converter func(*From) To) stream.Iter[To] {
	return collection.FilterAndConvert(elements, not.Nil[From], converter)
}
