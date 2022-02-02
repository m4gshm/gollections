package K

import "github.com/m4gshm/gollections/map_"

//V is K.V shortening for map_.NewKV
func V[k comparable, v any](key k, value v) *map_.KV[k, v] {
	return map_.NewKV(key, value)
}
