// Package conv provides slice converation helpers
package conv

import (
	"github.com/m4gshm/gollections/slice"
)

// AndReduce - conv.AndReduce converts elements and merge them into one
func AndReduce[FS ~[]From, From, To any](elements FS, converter func(From) (To, error), merge func(To, To) To) (To, error) {
	return slice.ConvAndReduce(elements, converter, merge)
}
