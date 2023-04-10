package not

import "github.com/m4gshm/gollections/check"

// not.Nil checks a reference for no nil value.
func Nil[T any](val *T) bool {
	return check.NotNil(val)
}