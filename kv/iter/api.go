// Package iter provides generic constructors and helpers for key/value iterators
package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/kv"
	kvLoop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/loop/group"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/slice"
)

// OfPairs instantiates KVIterator of predefined key\value pairs
func OfPairs[K, V any](pairs ...c.KV[K, V]) loop.KeyValuer[c.KV[K, V], K, V] {
	return WrapPairs(pairs)
}

// WrapPairs instantiates KVIterator using slice as the key\value pairs source
func WrapPairs[K, V any, P ~[]c.KV[K, V]](pairs P) loop.KeyValuer[c.KV[K, V], K, V] {
	return FromPairs[K, V](slice.NewIter(pairs))
}

// FromPairs converts an iterator of key\value pair elements to a KVIterator
func FromPairs[K, V any, I c.Iterator[c.KV[K, V]]](elements I) loop.KeyValuer[c.KV[K, V], K, V] {
	return FromIter(elements, (c.KV[K, V]).Key, (c.KV[K, V]).Value)
}

// FromIter converts a c.Iterator to a kv.KVIterator using key and value extractors
func FromIter[T, K, V any, I c.Iterator[T]](elements I, keyExtractor func(T) K, valExtractor func(T) V) loop.KeyValuer[T, K, V] {
	return iter.ToKV(elements, keyExtractor, valExtractor)
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any, I kv.KVIterator[K, V]](elements I) map[K][]V {
	return group.Of(elements.Next)
}

// Map instantiates key/value iterator that converts elements with a converter and returns them
func Map[K comparable, V any, KOUT comparable, VOUT any, I kv.KVIterator[K, V]](elements I, by func(K, V) (KOUT, VOUT)) stream.Iter[KOUT, VOUT, map[KOUT]VOUT] {
	return stream.New(kvLoop.Convert(elements.Next, by).Next, kvLoop.ToMap[KOUT, VOUT])
}

// Filter instantiates key/value iterator that iterates only over filtered elements
func Filter[K comparable, V any, I kv.KVIterator[K, V]](elements I, filter func(K, V) bool) stream.Iter[K, V, map[K]V] {
	return stream.New(kvLoop.Filter(elements.Next, filter).Next, kvLoop.ToMap[K, V])
}

// FilterKey instantiates key/value iterator that iterates only over elements that filtered by the key
func FilterKey[K comparable, V any, I kv.KVIterator[K, V]](elements I, fit func(K) bool) stream.Iter[K, V, map[K]V] {
	return Filter(elements, filter.Key[V](fit))
}

// FilterValue instantiates key/value iterator that iterates only over elements that filtered by the value
func FilterValue[K comparable, V any, I kv.KVIterator[K, V]](elements I, fit func(V) bool) stream.Iter[K, V, map[K]V] {
	return Filter(elements, filter.Value[K](fit))
}

// Reduce reduces keys/value pairs to an one pair
func Reduce[K comparable, V any, I kv.KVIterator[K, V]](elements I, by c.Quaternary[K, V]) (K, V) {
	return kvLoop.Reduce(elements.Next, by)
}
