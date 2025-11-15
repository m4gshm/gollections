// Package desc provides aliases for storing slice functions.
package desc

import (
	"cmp"

	"github.com/m4gshm/gollections/slice/sort"
)

// Of is a short alias for sort.DescBy
func Of[T any, O cmp.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return sort.DescBy(elements, orderConverner)
}
