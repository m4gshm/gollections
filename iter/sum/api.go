// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
)

// Of an alias of the it.Sum
func Of[T c.Summable, IT c.Iterator[T]](elements IT) T {
	return loop.Sum(elements.Next)
}
