// Package sum provides sum.Of alias
package sum

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/slice/sum"
)

func Of[T c.Summable](elements ...T) func() T {
	return func() T { return sum.Of(elements) }
}
