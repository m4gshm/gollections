// Package first provides short aliases for loop functions for retrieving a first element
package first

import (
	"github.com/m4gshm/gollections/loop"
)

// Of an alias of the loop.First
func Of[T any](next func() (T, bool), predicate func(T) bool) (T, bool) {
	return loop.First(next, predicate)
}

// Converted converts the first element that satisfies the condition of the 'predicate' function by the converter and returns it
func Converted[From, To any](next func() (From, bool), filter func(From) bool, converter func(From) To) (out To, ok bool) {
	if f, ok := loop.First(next, filter); ok {
		return converter(f), true
	}
	return out, false
}
