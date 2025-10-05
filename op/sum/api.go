// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice/sum"
)

// Of an alias of the slice.Sum
func Of[T op.Summable](elements ...T) T {
	return sum.Of(elements)
}
