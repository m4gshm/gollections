package clone

import (
	"github.com/m4gshm/gollections/slice"
)

// Of - synonym of the slice.Clone
func Of[TS ~[]T, T any](elements TS) TS {
	return slice.Clone(elements)
}

// Deep - synonym of the slice.DeepClone
func Deep[TS ~[]T, T any](elements TS, copier func(T) T) TS {
	return slice.DeepClone(elements, copier)
}
