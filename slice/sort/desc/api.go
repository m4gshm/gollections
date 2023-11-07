package desc

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice/sort"
)

func Of[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return sort.DescBy(elements, orderConverner)
}
