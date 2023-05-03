// Package kv provides key/value types, functions
package kv

import "github.com/m4gshm/gollections/c"

// Iterator provides iterate over key/value pairs, where an iteration can be interrupted by an error
type Iterator[K, V any] interface {
	// Next returns the next key/value pair.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (key K, value V, ok bool, err error)
	c.TrackLoop[K, V]
}
