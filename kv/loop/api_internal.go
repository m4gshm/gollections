package loop

import (
	"github.com/m4gshm/gollections/kv"
)

func startKvIt[K, V any, I kv.Iterator[K, V]](i I) (I, K, V, bool) {
	k, v, ok := i.Next()
	return i, k, v, ok
}
