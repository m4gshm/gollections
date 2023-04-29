package collection

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
)

// Vector - collection interface that provides elements order and access by index to the elements.
type Vector[T any] interface {
	c.Collection[T]

	c.TrackLoop[int, T]
	c.TrackEachLoop[int, T]

	c.Access[int, T]

	Len() int
	IsEmpty() bool
}

// Set - collection interface that ensures the uniqueness of elements (does not insert duplicate values).
type Set[T comparable] interface {
	c.Collection[T]
	c.Checkable[T]

	Len() int
	IsEmpty() bool
}

// Map - collection interface that stores key/value pairs and provide access to an element by its key
type Map[K comparable, V any] interface {
	kv.Collection[K, V, map[K]V]
	c.Checkable[K]
	c.Access[K, V]

	Len() int
	IsEmpty() bool
}
