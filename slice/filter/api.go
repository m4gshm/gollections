package filter

import (
	"github.com/m4gshm/gollections/slice"
)

// AndConvert filters the 'From' elements, and then converts them to 'To'
func AndConvert[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) []To {
	return slice.FilterAndConvert(elements, filter, converter)
}

// ConvertFilter filters the 'From' elements, then converts them to 'To', and then filters that ones
func ConvertFilter[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To, filterTo func(To) bool) []To {
	return slice.FilterConvertFilter(elements, filter, converter, filterTo)
}
