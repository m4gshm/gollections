// Package filter provides aliases for loop filtering helpers
package filter

import (
	"github.com/m4gshm/gollections/loop"
)

// AndConvert filters the 'From' elements, and then converts them to 'To'
func AndConvert[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To) loop.Loop[To] {
	return loop.FilterAndConvert(next, filter, converter)
}

// ConvertFilter filters the 'From' elements, then converts them to 'To', and then filters that ones
func ConvertFilter[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To, filterTo func(To) bool) loop.Loop[To] {
	return loop.FilterConvertFilter(next, filter, converter, filterTo)
}
