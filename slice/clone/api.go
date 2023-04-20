package clone

import (
	"github.com/m4gshm/gollections/ptr"
	"github.com/m4gshm/gollections/slice"
)

// Of is synonym of the slice.Clone
func Of[TS ~[]T, T any](elements TS) TS {
	return slice.Clone(elements)
}

// Deep is synonym of the slice.DeepClone
func Deep[TS ~[]T, T any](elements TS, copier func(T) T) TS {
	return slice.DeepClone(elements, copier)
}

// Ptr returns a pointer to a copy of the value pointed to by 'p'
func Ptr[T any](p *T) *T { return ptr.Of(*p) }
