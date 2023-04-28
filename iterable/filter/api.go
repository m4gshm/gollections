// Package filter provides aliases for collections filtering helpers
package filter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iterable"
	"github.com/m4gshm/gollections/loop"
)

// AndConvert - filter.AndConvert is short alias of iterable.FilterAndConvert
func AndConvert[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, filter func(From) bool, converter func(From) To) loop.StreamIter[To] {
	return iterable.FilterAndConvert(elements, filter, converter)
}

// ConvertFilter - filter.ConvertFilter is short alias of slice.FilterConvertFilter
func ConvertFilter[From, To any, IT c.Iterable[c.Iterator[From]]](elements IT, filterFrom func(From) bool, converter func(From) To, filterTo func(To) bool) loop.StreamIter[To] {
	return iterable.FilterAndConvert(elements, filterFrom, converter).Filter(filterTo)
}
