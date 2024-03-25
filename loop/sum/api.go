// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// Of an alias of the loop.Sum
func Of[T c.Summable](sum func() (T, bool)) T {
	return loop.Sum(sum)
}
