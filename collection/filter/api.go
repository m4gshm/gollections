// Package filter provides aliases for collections filtering helpers
package filter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/stream"
)

// AndConvert - filter.AndConvert is short alias of collection.FilterAndConvert
func AndConvert[From, To any, IT c.Iterable[From]](elements IT, filter func(From) bool, converter func(From) To) stream.Iter[To] {
	return collection.FilterAndConvert(elements, filter, converter)
}

// ConvertFilter - filter.ConvertFilter is short alias of slice.FilterConvertFilter
func ConvertFilter[From, To any, IT c.Iterable[From]](elements IT, filterFrom func(From) bool, converter func(From) To, filterTo func(To) bool) stream.Iter[To] {
	return collection.FilterAndConvert(elements, filterFrom, converter).Filter(filterTo)
}
