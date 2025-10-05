package collection

import (
	"github.com/m4gshm/gollections/c"
)

// Iterator provides iterate over key/value pairs
type Iterator[K, V any] interface {
	// Next returns the next key/value pair.
	// The ok result indicates whether the element was returned by the iterator.
	// If ok == false, then the iteration must be completed.
	Next() (key K, value V, ok bool)
	c.TrackEach[K, V]
}

// Collection is the base interface of associative collections
type Collection[K comparable, V any, M map[K]V | map[K][]V] interface {
	c.TrackEach[K, V]
	c.MapFactory[K, V, M]

	Reduce(merge func(K, K, V, V) (K, V)) (K, V)
	HasAny(func(K, V) bool) bool
	All(consumer func(K, V) bool)
}

// Convertable provides limited kit of map transformation methods
type Convertable[K, V any,
	Seq2 ~func(yield func(K, V) bool),
	SeqE ~func(yield func(c.KV[K, V], error) bool),
] interface {
	Convert(converter func(K, V) (K, V)) Seq2
	Conv(converter func(K, V) (K, V, error)) SeqE

	ConvertKey(converter func(K) K) Seq2
	ConvertValue(converter func(V) V) Seq2

	ConvKey(converter func(K) (K, error)) SeqE
	ConvValue(converter func(V) (V, error)) SeqE
}

// Filterable provides limited kit of filering methods
type Filterable[K, V any,
	Seq2 ~func(yield func(K, V) bool),
	SeqE ~func(yield func(c.KV[K, V], error) bool),
] interface {
	Filter(predicate func(K, V) bool) Seq2
	Filt(predicate func(K, V) (bool, error)) SeqE

	FilterKey(predicate func(K) bool) Seq2
	FilterValue(predicate func(V) bool) Seq2

	FiltKey(predicate func(K) (bool, error)) SeqE
	FiltValue(predicate func(V) (bool, error)) SeqE
}
