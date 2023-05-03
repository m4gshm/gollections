// Package range_ provides alias for the slice.Range function
package range_

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice"
)

// Of - range_.Of short alias of the slice.Range
func Of[T constraints.Integer](from T, to T) []T {
	return slice.Range(from, to)
}
