// Package range_ provides alias for the slice.Range function
package range_

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/loop"
)

// Of - range_.Of short alias of the loop.Range
func Of[T constraints.Integer](from T, to T) func() (T, bool) {
	return loop.Range(from, to)
}
