// Package filter provides aliases for collections filtering helpers
package filter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/stream"
)

// AndConvert - filter.AndConvert is short alias of iterable.FilterAndConvert
func AndConvert[I c.Iterator[From],From, To any, IT c.Iterable[I]](elements IT, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	return iterable.FilterAndConvert[I](elements, filter, converter)
}

// ConvertFilter - filter.ConvertFilter is short alias of slice.FilterConvertFilter
func ConvertFilter[I c.Iterator[From], From, To any, IT c.Iterable[I]](elements IT, filterFrom func(From) bool, converter func(From) To, filterTo func(To) bool) stream.Iter[To] {
	return iterable.FilterAndConvert[I](elements, filterFrom, converter).Filter(filterTo)
}
