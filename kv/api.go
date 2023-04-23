package kv

import "github.com/m4gshm/gollections/c"

// New creates a key/value pair holder
func New[K any, V any](key K, value V) c.KV[K, V] {
	return c.KV[K, V]{K: key, V: value}
}
