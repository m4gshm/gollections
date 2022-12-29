package as

import (
	"github.com/m4gshm/gollections/conv"
)

// Is an alias of the conv.AsIs
func Is[T any](value T) T { return conv.AsIs(value) }
