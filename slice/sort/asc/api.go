package asc

import (
	"github.com/m4gshm/gollections/slice/sort"
	"golang.org/x/exp/constraints"
)

func Of[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return sort.By(elements, orderConverner)
}
