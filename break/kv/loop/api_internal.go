package loop

import (
	"github.com/m4gshm/gollections/break/kv"
)

func startKvIt[K, V any, I kv.Iterator[K, V]](i I) (I, K, V, bool, error) {
	k, v, ok, err := i.Next()
	return i, k, v, ok, err
}
