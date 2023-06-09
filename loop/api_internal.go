package loop

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
)

func startIt[T any, I c.Iterator[T]](i I) (I, T, bool) {
	element, ok := i.Next()
	return i, element, ok
}

func startKvIt[K, V any, I kv.Iterator[K, V]](i I) (I, K, V, bool) {
	k, v, ok := i.Next()
	return i, k, v, ok
}
