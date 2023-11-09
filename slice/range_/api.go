// Package range_ provides alias for the slice.Range function
package range_

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice"
)

// Of - range_.Of short alias of the slice.Range
func Of[T constraints.Integer | rune](from T, to T) []T {
	return slice.Range(from, to)
}

// Closed - range_.Closed short alias of the slice.RangeClosed
func Closed[T constraints.Integer | rune](from T, to T) []T {
	return slice.RangeClosed(from, to)
}
