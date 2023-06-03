// Package val provides pointer to value convert helpers
package val

import "github.com/m4gshm/gollections/convert"

// Of is pointer-tovalue conversion helper
func Of[T any](t *T) T {
	return convert.PtrVal(t)
}
