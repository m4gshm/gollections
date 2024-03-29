// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice"
)

// Of an alias of the slice.Sum
func Of[TS ~[]T, T c.Summable](elements TS) T {
	return slice.Sum(elements)
}
