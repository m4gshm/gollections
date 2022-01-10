package range_

import (
	"constraints"

	"github.com/m4gshm/gollections/slice"
)

func Of[T constraints.Integer](from T, to T) []T {
	return slice.Range(from, to)
}
