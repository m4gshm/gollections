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

// IterFor extends an iterator type by a 'Start' function implementation
type IterFor[K, V any, I Iterator[K, V]] interface {
	// Start is used with for loop construct.
	// Returns the iterator itself, the first key/value pair, ok == false if the iteration must be completed, and an error.
	//
	// 	var i IterFor = ...
	//	for i, k, v, ok, err := i.Start(); ok || err != nil; k, v, ok, err = i.Next() {
	//		if err != nil {
	//			return err
	//		}
	//	}
	Start() (iterator I, key K, value V, ok bool, err error)
}
