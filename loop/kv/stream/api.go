package stream

import (
	"github.com/m4gshm/gollections/as"
	"github.com/m4gshm/gollections/kviter/group"
	"github.com/m4gshm/gollections/loop/iter"
)

// Group transforms iterable elements to the MapPipe based on applying key extractor to the elements
func Group[T any, K comparable](next func() (T, bool), keyExtractor func(T) K) KVStream[K, T, map[K][]T] {
	return GroupAndConvert(next, keyExtractor, as.Is[T])
}

// GroupAndConvert transforms iterable elements to the MapPipe based on applying key extractor to the elements
func GroupAndConvert[T any, K comparable, V any](next func() (T, bool), keyExtractor func(T) K, valueConverter func(T) V) KVStream[K, V, map[K][]V] {
	kv := iter.NewKeyValuer(next, keyExtractor, valueConverter)
	return New(kv.Next, group.Of[K, V])
}
