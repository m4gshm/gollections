// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/sum"
)

// Of an alias of the slice.Sum
func Of[T c.Summable](elements ...T) T {
	return sum.Of(elements)
}
