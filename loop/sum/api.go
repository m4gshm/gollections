// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/op"
)

// Of an alias of the loop.Sum
func Of[T op.Summable](sum func() (T, bool)) T {
	return loop.Sum(sum)
}
