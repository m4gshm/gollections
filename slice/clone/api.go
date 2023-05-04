// Package clone provides slice clone aliases
package clone

import (
	"github.com/m4gshm/gollections/convert/ptr"
	"github.com/m4gshm/gollections/slice"
)

// Of is slice.Clone alias
func Of[TS ~[]T, T any](elements TS) TS {
	return slice.Clone(elements)
}

// Deep is slice.DeepClone alias
func Deep[TS ~[]T, T any](elements TS, copier func(T) T) TS {
	return slice.DeepClone(elements, copier)
}

// Ptr returns a pointer to a copy of the value pointed to by the 'p'
func Ptr[T any](p *T) *T { return ptr.Of(*p) }
