package kviter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter"
	implit "github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/kviter/group"
	"github.com/m4gshm/gollections/kviter/impl/kvit"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/op"
)

// OfPairs instantiates KVIterator of predefined key\value pairs
func OfPairs[K, V any](pairs ...c.KV[K, V]) c.KVIterator[K, V] {
	return WrapPairs(pairs)
}

// WrapPairs instantiates KVIterator using slice as the key\value pairs source
func WrapPairs[K, V any, P ~[]c.KV[K, V]](pairs P) c.KVIterator[K, V] {
	return FromPairs[K, V](iter.Wrap(pairs))
}

// FromPairs converts a iterator of key\value pair elements to a KVIterator
func FromPairs[K, V any](elements c.Iterator[c.KV[K, V]]) c.KVIterator[K, V] {
	return FromIter(elements, (c.KV[K, V]).Key, (c.KV[K, V]).Value)
}

// FromIter converts a c.Iterator to a c.KVIterator using key and value extractors
func FromIter[T, K, V any](elements c.Iterator[T], keyExtractor func(T) K, valExtractor func(T) V) c.KVIterator[K, V] {
	return iter.ToPairs(elements, keyExtractor, valExtractor)
}

// FirstVal - ToMap value resolver
func FirstVal[K, V any](exists bool, key K, old, new V) V { return op.IfElse(exists, old, new) }

// LastVal - ToMap value resolver
func LastVal[K, V any](exists bool, key K, old, new V) V { return new }

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[K comparable, V any](it c.KVIterator[K, V]) map[K]V {
	return ToMapResolv(it, FirstVal[K, V])
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[K comparable, E, V any](it c.KVIterator[K, E], valResolv func(bool, K, V, E) V) map[K]V {
	e := map[K]V{}
	for k, elem, ok := it.Next(); ok; k, elem, ok = it.Next() {
		exists, ok := e[k]
		e[k] = valResolv(ok, k, exists, elem)
	}
	return e
}

// Group collects sets of values grouped by keys obtained by passing a key/value iterator
func Group[K comparable, V any](it c.KVIterator[K, V]) map[K][]V {
	return group.Of(it)
}

// OfLoop creates an IteratorBreakable instance that loops over elements of a source
// The hasNext specifies a predicate that tests existing of a next element in the source.
// The getNext extracts the element.
func OfLoop[S, k, V any](source S, hasNext func(S) bool, getNext func(S) (k, V, error)) c.KVIteratorBreakable[k, V] {
	l := kvit.NewLoop(source, hasNext, getNext)
	return &l
}

// Map instantiates key/value iterator that converts elements with a converter and returns them
func Map[K comparable, V any, Kto comparable, Vto any](elements c.KVIterator[K, V], by func(K, V) (Kto, Vto)) c.MapPipe[Kto, Vto, map[Kto]Vto] {
	return implit.NewKVPipe(implit.ConvertKV(elements, by), ToMap[Kto, Vto])
}

// Filter instantiates key/value iterator that iterates only over filtered elements
func Filter[K comparable, V any, IT c.KVIterator[K, V]](elements IT, filter func(K, V) bool) c.MapPipe[K, V, map[K]V] {
	return implit.NewKVPipe(implit.FilterKV(elements, filter), ToMap[K, V])
}

// FilterKey instantiates key/value iterator that iterates only over elements that filtered by the key
func FilterKey[K comparable, V any](elements c.KVIterator[K, V], fit func(K) bool) c.MapPipe[K, V, map[K]V] {
	return Filter(elements, filter.Key[V](fit))
}

// FilterValue instantiates key/value iterator that iterates only over elements that filtered by the value
func FilterValue[K comparable, V any](elements c.KVIterator[K, V], fit func(V) bool) c.MapPipe[K, V, map[K]V] {
	return Filter(elements, filter.Value[K](fit))
}

// Reduce reduces keys/value pairs to an one pair
func Reduce[K comparable, V any](elements c.KVIterator[K, V], by c.Quaternary[K, V]) (K, V) {
	return loop.ReduceKV(elements.Next, by)
}
