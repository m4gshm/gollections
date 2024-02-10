// Package kv provides generic key/value pair constructors and helpers
package kv

import "github.com/m4gshm/gollections/c"

// New creates a key/value pair holder
func New[K any, V any](key K, value V) c.KV[K, V] {
	return c.KV[K, V]{K: key, V: value}
}

func All[K, V any](next func() (K, V, bool), yield func(K, V) bool) {
	for k, v, ok := next(); ok && yield(k, v); k, v, ok = next() {
	}
}
