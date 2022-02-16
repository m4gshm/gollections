package K

import (
	"github.com/m4gshm/gollections/c"
)

//V is K.V shortening for map_.NewKV
func V[k comparable, v any](key k, value v) *c.KV[k, v] {
	return c.NewKV(key, value)
}
