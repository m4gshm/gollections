// Package c provides common types and functions
package c

import "github.com/m4gshm/gollections/c"

// Iterator provides iterate over elements of a source, where an iteration can be interrupted by an error
type Iterator[T any] interface {
	// Next returns the next element.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (out T, ok bool, err error)

	c.ForLoop[T]
}

// IterFor extends an iterator type by a 'Start' function implementation
type IterFor[T any, I Iterator[T]] interface {
	// Start is used with for loop construct.
	// Returns the iterator itself, the first element, ok == false if the iteration must be completed, and an error.
	//
	// 	var i IterFor = ...
	//	for i, val, ok, err := i.Start(); ok || err != nil; val, ok, err = i.Next() {
	//		if err != nil {
	//			return err
	//		}
	//	}
	Start() (iterator I, val T, ok bool, err error)
}
