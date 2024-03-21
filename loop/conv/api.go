// Package conv provides loop converation helpers
package conv

import (
	breakLoop "github.com/m4gshm/gollections/break/loop"
	"github.com/m4gshm/gollections/loop"
)

// FromIndexed - conv.FromIndexed retrieves elements from a indexed source and converts them
func FromIndexed[From, To any](len int, next func(int) From, converter func(from From) (To, error)) breakLoop.Loop[To] {
	return loop.Conv(loop.OfIndexed(len, next), converter)
}

// AndReduce - convert.AndReduce converts elements and merge them into one
func AndReduce[From, To any](next func() (From, bool), converter func(From) (To, error), merge func(To, To) To) (To, error) {
	return loop.ConvAndReduce(next, converter, merge)
}
