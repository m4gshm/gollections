// Package as provides as.Is alias
package as

import (
	"github.com/m4gshm/gollections/convert"
)

// Is an alias of the conv.AsIs
func Is[T any](value T) T { return convert.AsIs(value) }

// Is an alias of the conv.AsSlice
func Slice[T any](value T) []T { return convert.AsSlice(value) }
