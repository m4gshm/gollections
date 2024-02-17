// Package range_ provides alias for the slice.Range function
package range_

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/loop"
)

// Of - range_.Of short alias of the loop.Range
func Of[T constraints.Integer](from T, to T) loop.Loop[T] {
	return loop.Range(from, to)
}

// Closed - range_.Closed short alias of the loop.RangeClosed
func Closed[T constraints.Integer](from T, to T) loop.Loop[T] {
	return loop.RangeClosed(from, to)
}
