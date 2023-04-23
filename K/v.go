package k

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
)

// V is k.V shortening for map_.NewKV
func V[k, v any](key k, value v) c.KV[k, v] {
	return kv.New(key, value)
}
