// Package sort provides sorting of cloned slice elements
package sort

import (
	"golang.org/x/exp/constraints"

	"github.com/m4gshm/gollections/slice/clone"
	"github.com/m4gshm/gollections/slice/sort"
)

// By sorts cloned elements slice in ascending order, using the orderConverner function to retrieve a value of type Ordered.
func By[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return sort.By(clone.Of(elements), orderConverner)
}

// DescBy sorts cloned elements slice in descending order, using the orderConverner function to retrieve a value of type Ordered.
func DescBy[T any, O constraints.Ordered, TS ~[]T](elements TS, orderConverner func(T) O) TS {
	return sort.DescBy(clone.Of(elements), orderConverner)
}

// Asc sorts orderable elements ascending
func Asc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return sort.Asc(clone.Of(elements))
}

// Desc sorts orderable elements descending
func Desc[T constraints.Ordered, TS ~[]T](elements TS) TS {
	return sort.Desc(clone.Of(elements))
}
