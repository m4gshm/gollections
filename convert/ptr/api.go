// Package ptr provides value, pointer convert helpers
package ptr

import "github.com/m4gshm/gollections/convert"

// Of is value-to-pointer conversion helper
func Of[T any](t T) *T {
	return convert.ToPointer(t)
}
