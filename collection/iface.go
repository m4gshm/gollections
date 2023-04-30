package collection

import (
	breakKvstream "github.com/m4gshm/gollections/break/kv/stream"
	breakStream "github.com/m4gshm/gollections/break/stream"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv"
	kvstream "github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/stream"
)

type Collection[T any] interface {
	c.Collection[T]
	c.Filterable[T, stream.Iter[T], breakStream.Iter[T]]
	c.Convertable[T, stream.Iter[T], breakStream.Iter[T]]

	Len() int
	IsEmpty() bool

	HasAny(func(T) bool) bool
}

// Vector - collection interface that provides elements order and access by index to the elements.
type Vector[T any] interface {
	Collection[T]

	c.TrackLoop[int, T]
	c.TrackEachLoop[int, T]

	c.Access[int, T]
}

// Set - collection interface that ensures the uniqueness of elements (does not insert duplicate values).
type Set[T comparable] interface {
	Collection[T]
	c.Checkable[T]
}

// Map - collection interface that stores key/value pairs and provide access to an element by its key
type Map[K comparable, V any] interface {
	kv.Collection[K, V, map[K]V]
	kv.Filterable[K, V, kvstream.Iter[K, V, map[K]V], breakKvstream.Iter[K, V, map[K]V]]
	kv.Convertable[K, V, kvstream.Iter[K, V, map[K]V], breakKvstream.Iter[K, V, map[K]V]]
	c.Checkable[K]
	c.Access[K, V]

	Len() int
	IsEmpty() bool

	HasAny(func(K, V) bool) bool
}
