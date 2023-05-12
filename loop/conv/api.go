// Package conv provides loop converation helpers
package conv

import (
	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/loop"
)

// FromIndexed - conv.FromIndexed retrieves elements from a indexed source and converts them
func FromIndexed[From, To any](len int, next func(int) From, converter func(from From) (To, error)) breakLoop.ConvertIter[From, To] {
	return loop.Conv(loop.OfIndexed(len, next), converter)
}
