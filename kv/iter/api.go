// Package iter provides generic constructors and helpers for key/value iterators
package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
	"github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/loop/group"
	"github.com/m4gshm/gollections/map_/filter"
)

// OfPairs instantiates KVIterator of predefined key\value pairs
func OfPairs[K, V any](pairs ...c.KV[K, V]) c.KVIterator[K, V] {
	return WrapPairs(pairs)
}

// WrapPairs instantiates KVIterator using slice as the key\value pairs source
func WrapPairs[K, V any, P ~[]c.KV[K, V]](pairs P) c.KVIterator[K, V] {
	return FromPairs[K, V](iter.Wrap(pairs))
}

// FromPairs converts an iterator of key\value pair elements to a KVIterator
func FromPairs[K, V any](elements c.Iterator[c.KV[K, V]]) c.KVIterator[K, V] {
	return FromIter(elements, (c.KV[K, V]).Key, (c.KV[K, V]).Value)
}

// FromIter converts a c.Iterator to a c.KVIterator using key and value extractors
func FromIter[T, K, V any](elements c.Iterator[T], keyExtractor func(T) K, valExtractor func(T) V) c.KVIterator[K, V] {
	return iter.ToPairs(elements, keyExtractor, valExtractor)
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	return group.Of(it.Next)
}

// OfLoop creates an IteratorBreakable instance that loops over elements of a source
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, K, V any](source S, hasNext func(S) bool, getNext func(S) (K, V, error)) c.KVIteratorBreakable[K, V] {
	l := loop.NewIter(source, hasNext, getNext)
	return &l
}

// Map instantiates key/value iterator that converts elements with a converter and returns them
func Map[K comparable, V any, KOUT comparable, VOUT any](elements c.KVIterator[K, V], by func(K, V) (KOUT, VOUT)) c.KVStream[KOUT, VOUT, map[KOUT]VOUT] {
	return loop.Stream(loop.Convert(elements.Next, by).Next, loop.ToMap[KOUT, VOUT])
}

// Filter instantiates key/value iterator that iterates only over filtered elements
func Filter[K comparable, V any, IT c.KVIterator[K, V]](elements IT, filter func(K, V) bool) c.KVStream[K, V, map[K]V] {
	return loop.Stream(loop.Filter(elements.Next, filter).Next, loop.ToMap[K, V])
}

// FilterKey instantiates key/value iterator that iterates only over elements that filtered by the key
func FilterKey[K comparable, V any](elements c.KVIterator[K, V], fit func(K) bool) c.KVStream[K, V, map[K]V] {
	return Filter(elements, filter.Key[V](fit))
}

// FilterValue instantiates key/value iterator that iterates only over elements that filtered by the value
func FilterValue[K comparable, V any](elements c.KVIterator[K, V], fit func(V) bool) c.KVStream[K, V, map[K]V] {
	return Filter(elements, filter.Value[K](fit))
}

// Reduce reduces keys/value pairs to an one pair
func Reduce[K comparable, V any](elements c.KVIterator[K, V], by c.Quaternary[K, V]) (K, V) {
	return loop.Reduce(elements.Next, by)
}
