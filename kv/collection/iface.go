package collection

import (
	breakloop "github.com/m4gshm/gollections/break/kv/loop"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/kv/loop"
)

// Iterator provides iterate over key/value pairs
type Iterator[K, V any] interface {
	// Next returns the next key/value pair.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (key K, value V, ok bool)
	c.Track[K, V]
	c.TrackEach[K, V]
}

// Iterable is an iterator supplier interface
type Iterable[K, V any] interface {
	Loop() loop.Loop[K, V]
}

// Collection is the base interface of associative collections
type Collection[K comparable, V any, M map[K]V | map[K][]V] interface {
	c.Track[K, V]
	c.TrackEach[K, V]
	Iterable[K, V]
	c.MapFactory[K, V, M]

	Reduce(merge func(K, K, V, V) (K, V)) (K, V)
	HasAny(func(K, V) bool) bool
	All(consumer func(K, V) bool)
}

// Convertable provides limited kit of map transformation methods
type Convertable[K, V any] interface {
	Convert(converter func(K, V) (K, V)) loop.Loop[K, V]
	Conv(converter func(K, V) (K, V, error)) breakloop.Loop[K, V]

	ConvertKey(converter func(K) K) loop.Loop[K, V]
	ConvertValue(converter func(V) V) loop.Loop[K, V]

	ConvKey(converter func(K) (K, error)) breakloop.Loop[K, V]
	ConvValue(converter func(V) (V, error)) breakloop.Loop[K, V]
}

// Filterable provides limited kit of filering methods
type Filterable[K, V any] interface {
	Filter(predicate func(K, V) bool) loop.Loop[K, V]
	Filt(predicate func(K, V) (bool, error)) breakloop.Loop[K, V]

	FilterKey(predicate func(K) bool) loop.Loop[K, V]
	FilterValue(predicate func(V) bool) loop.Loop[K, V]

	FiltKey(predicate func(K) (bool, error)) breakloop.Loop[K, V]
	FiltValue(predicate func(V) (bool, error)) breakloop.Loop[K, V]
}
