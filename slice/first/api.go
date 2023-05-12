// Package first provides short aliases for slice functions for retrieving a first element
package first

import (
	"github.com/m4gshm/gollections/slice"
)

// Of returns the first element that satisfies requirements of the specified filter
func Of[TS ~[]T, T any](elements TS, filter func(T) bool) (T, bool) {
	return slice.First(elements, filter)
}

// Converted converts the first element that satisfies the condition of the 'filter' function by the converter and returns it
func Converted[FS ~[]From, From, To any](elements FS, filter func(From) bool, converter func(From) To) (out To, ok bool) {
	if f, ok := slice.First(elements, filter); ok {
		return converter(f), true
	}
	return out, false
}
