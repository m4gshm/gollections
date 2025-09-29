package collection

import (
	"github.com/m4gshm/gollections/c"
	kv "github.com/m4gshm/gollections/kv/collection"
)

// Collection is the base interface for the Vector and the Set impelementations
type Collection[T any] interface {
	c.Collection[T]
	c.Filterable[T, Seq[T], SeqE[T]]
	c.Convertable[T, Seq[T], SeqE[T]]

	Len() int
	IsEmpty() bool

	HasAny(func(T) bool) bool
}

// Vector - collection interface that provides elements order and access by index to the elements.
type Vector[T any] interface {
	Collection[T]

	c.Track[int, T]
	c.TrackEach[int, T]

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
	kv.Filterable[K, V, Seq2[K, V], SeqE[c.KV[K, V]]]
	kv.Convertable[K, V, Seq2[K, V], SeqE[c.KV[K, V]]]
	c.Checkable[K]
	c.Access[K, V]
	c.KVRange[K, V]

	Len() int
	IsEmpty() bool

	HasAny(func(K, V) bool) bool
}
