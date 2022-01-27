package K

import "github.com/m4gshm/gollections/typ"

//V is K.V shortening for typ.NewKV
func V[k comparable, v any](key k, value v) *typ.KV[k, v] {
	return typ.NewKV(key, value)
}
