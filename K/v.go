package K

import "github.com/m4gshm/container/typ"

func V[k comparable, v any](key k, value v) *typ.KV[k, v] {
	return typ.NewKV(key, value)
}
