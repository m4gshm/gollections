package loop

import (
	"github.com/m4gshm/gollections/break/c"
	"github.com/m4gshm/gollections/break/kv"
)

func startIt[T any, I c.Iterator[T]](i I) (I, T, bool, error) {
	element, ok, err := i.Next()
	return i, element, ok, err
}

func startKvIt[K, V any, I kv.Iterator[K, V]](i I) (I, K, V, bool, error) {
	k, v, ok, err := i.Next()
	return i, k, v, ok, err
}
