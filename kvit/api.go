package kvit

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it"
	implit "github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/kvit/group"
	"github.com/m4gshm/gollections/kvit/impl/kvit"
	"github.com/m4gshm/gollections/ptr"
)

// OfPairs instantiates KVIterator of predefined key\value pairs
func OfPairs[K, V any](pairs ...c.KV[K, V]) c.KVIterator[K, V] {
	return WrapPairs(pairs)
}

// WrapPairs instantiates KVIterator using slice as the key\value pairs source
func WrapPairs[K, V any, P ~[]c.KV[K, V]](pairs P) c.KVIterator[K, V] {
	return FromPairs(it.Wrap(pairs))
}

func FromPairs[K, V any](elements c.Iterator[c.KV[K, V]]) c.KVIterator[K, V] {
	return FromIter(elements, (c.KV[K, V]).Key, (c.KV[K, V]).Value)
}

func FromIter[T, K, V any](elements c.Iterator[T], keyExtractor c.Converter[T, K], valExtractor c.Converter[T, V]) c.KVIterator[K, V] {
	return it.ToPairs(elements, keyExtractor, valExtractor)
}

// FirstVal - ToMap value resolver
func FirstVal[K, V any](key K, exist, new V) V { return exist }

// LastVal - ToMap value resolver
func LastVal[K, V any](key K, exist, new V) V { return new }

// ToMap collects key\value elements to a map by iterating over the elements
func ToMap[K comparable, V any](it c.KVIterator[K, V]) map[K]V {
	return ToMapResolv(it, FirstVal[K, V])
}

// ToMapResolv collects key\value elements to a map by iterating over the elements with resolving of duplicated key values
func ToMapResolv[K comparable, V any](it c.KVIterator[K, V], valResolv func(K, V, V) V) map[K]V {
	e := map[K]V{}
	for k, v, ok := it.Next(); ok; k, v, ok = it.Next() {
		if exists, ok := e[k]; ok {
			e[k] = valResolv(k, exists, v)
		} else {
			e[k] = v
		}
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
	return ptr.Of(kvit.NewLoop(source, hasNext, getNext))
}

// Map instantiates key/value iterator that converts elements with a converter and returns them
func Map[K comparable, V any, Kto comparable, Vto any](elements c.KVIterator[K, V], by c.BiConverter[K, V, Kto, Vto]) c.MapPipe[Kto, Vto, map[Kto]Vto] {
	return implit.NewKVPipe(implit.MapKV(elements, by), ToMap[Kto, Vto])
}

// Filter instantiates key/value iterator that iterates only over filtered elements
func Filter[K comparable, V any, IT c.KVIterator[K, V]](elements IT, filter c.BiPredicate[K, V]) c.MapPipe[K, V, map[K]V] {
	return implit.NewKVPipe(implit.FilterKV(elements, filter), ToMap[K, V])
}

// FilterKey instantiates key/value iterator that iterates only over elements that filtered by the key
func FilterKey[K comparable, V any](elements c.KVIterator[K, V], filter c.Predicate[K]) c.MapPipe[K, V, map[K]V] {
	return Filter(elements, c.FitKey[K, V](filter))
}

// FilterValue instantiates key/value iterator that iterates only over elements that filtered by the value
func FilterValue[K comparable, V any](elements c.KVIterator[K, V], filter c.Predicate[V]) c.MapPipe[K, V, map[K]V] {
	return Filter(elements, c.FitValue[K](filter))
}

// Reduce reduces keys/value pairs to an one pair
func Reduce[K comparable, V any](elements c.KVIterator[K, V], by c.Quaternary[K, V]) (K, V) {
	return implit.ReduceKV(elements, by)
}
