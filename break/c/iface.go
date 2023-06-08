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

type IterFor[T any, I Iterator[T]] interface {
	Start() (iterator I, val T, ok bool, err error)
}
