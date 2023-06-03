package iter

import (
	breakc "github.com/m4gshm/gollections/break/c"
	breakkv "github.com/m4gshm/gollections/break/kv"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
)

func startBreakIt[T any, I breakc.Iterator[T]](i I) (I, T, bool, error) {
	element, ok, err := i.Next()
	return i, element, ok, err
}

func startIt[T any, I c.Iterator[T]](i I) (I, T, bool) {
	element, ok := i.Next()
	return i, element, ok
}

func startBreakKvIt[K, V any, I breakkv.Iterator[K, V]](i I) (I, K, V, bool, error) {
	k, v, ok, err := i.Next()
	return i, k, v, ok, err
}

func startKvIt[K, V any, I kv.Iterator[K, V]](i I) (I, K, V, bool) {
	k, v, ok := i.Next()
	return i, k, v, ok
}
