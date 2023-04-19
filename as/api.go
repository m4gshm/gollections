package as

import (
	"github.com/m4gshm/gollections/convert"
)

// Is an alias of the conv.AsIs
func Is[T any](value T) T { return convert.AsIs(value) }
