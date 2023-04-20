package filter

import (
	"github.com/m4gshm/gollections/slice"
)

// AndConvert - filter.AndConvert is short alias of slice.FilterAndConvert
func AndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) []To {
	return slice.FilterAndConvert(elements, filter, converter)
}

// ConvertFilter - filter.ConvertFilter is short alias of slice.FilterConvertFilter
func ConvertFilter[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To, filterTo func(To) bool) []To {
	return slice.FilterConvertFilter(elements, filter, converter, filterTo)
}
