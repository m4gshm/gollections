package kv

import "github.com/m4gshm/gollections/c"

// Iterator provides iterate over key/value pairs
type Iterator[K, V any] interface {
	// Next returns the next key/value pair.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (key K, value V, ok bool)
	c.TrackLoop[K, V]
	c.TrackEachLoop[K, V]

	All(func(K, V) bool)
}

// IterFor extends an iterator type by a 'Start' function implementation
type IterFor[K, V any, I Iterator[K, V]] interface {
	// Start is used with for loop construct.
	// Returns the iterator itself, the first key/value pair, and ok == false if the iteration must be completed.
	//
	// 	var i IterFor = ...
	// 	for i, k, v, ok := i.Start(); ok; k, v, ok = i.Next() {
	//  	_ = val
	//	}
	Start() (iterator I, key K, value V, ok bool)
}

// Iterable is an iterator supplier interface
type Iterable[K, V any] interface {
	Iter() Iterator[K, V]
}

// Collection is the base interface of associative collections
type Collection[K comparable, V any, M map[K]V | map[K][]V] interface {
	c.TrackLoop[K, V]
	c.TrackEachLoop[K, V]
	Iterable[K, V]
	c.MapFactory[K, V, M]

	Reduce(merger func(K, K, V, V) (K, V)) (K, V)

	All(func(K, V) bool)
}

// Convertable provides limited kit of map transformation methods
type Convertable[K, V, KVStream, KVStreamBreakable any] interface {
	Convert(converter func(K, V) (K, V)) KVStream
	Conv(converter func(K, V) (K, V, error)) KVStreamBreakable

	ConvertKey(converter func(K) K) KVStream
	ConvertValue(converter func(V) V) KVStream

	ConvKey(converter func(K) (K, error)) KVStreamBreakable
	ConvValue(converter func(V) (V, error)) KVStreamBreakable
}

// Filterable provides limited kit of filering methods
type Filterable[K, V, KVStream, KVStreamBreakable any] interface {
	Filter(predicate func(K, V) bool) KVStream
	Filt(predicate func(K, V) (bool, error)) KVStreamBreakable

	FilterKey(predicate func(K) bool) KVStream
	FilterValue(predicate func(V) bool) KVStream

	FiltKey(predicate func(K) (bool, error)) KVStreamBreakable
	FiltValue(predicate func(V) (bool, error)) KVStreamBreakable
}
