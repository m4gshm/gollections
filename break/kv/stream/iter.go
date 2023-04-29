package stream

import (
	"github.com/m4gshm/gollections/break/kv/loop"
	breakLoop "github.com/m4gshm/gollections/break/loop"
	breakFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/map_/filter"
)

func New[K comparable, V any, M map[K]V | map[K][]V](next func() (K, V, bool, error), collector MapCollector[K, V, M]) Iter[K, V, M] {
	return Iter[K, V, M]{next: next, collector: collector}
}

// Iter is the key/value Iterator based stream implementation.
type Iter[K comparable, V any, M map[K]V | map[K][]V] struct {
	next      func() (K, V, bool, error)
	collector MapCollector[K, V, M]
}

var (
	_ c.KVIteratorBreakable[string, any]                 = (*Iter[string, any, map[string]any])(nil)
	_ c.KVStreamBreakable[string, any, map[string]any]   = (*Iter[string, any, map[string]any])(nil)
	_ c.KVStreamBreakable[string, any, map[string][]any] = (*Iter[string, any, map[string][]any])(nil)

	_ c.KVIteratorBreakable[string, any]                 = Iter[string, any, map[string]any]{}
	_ c.KVStreamBreakable[string, any, map[string]any]   = Iter[string, any, map[string]any]{}
	_ c.KVStreamBreakable[string, any, map[string][]any] = Iter[string, any, map[string][]any]{}
)

// Next implements c.KVIterator
func (k Iter[K, V, M]) Next() (K, V, bool, error) {
	return k.next()
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (k Iter[K, V, M]) FilterKey(predicate func(K) bool) Iter[K, V, M] {
	return New(loop.Filter(k.next, filter.Key[V](predicate)).Next, k.collector)
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (k Iter[K, V, M]) FiltKey(predicate func(K) (bool, error)) Iter[K, V, M] {
	return New(loop.Filt(k.next, breakFilter.Key[V](predicate)).Next, k.collector)
}

// // ConvertKey returns a stream that applies the 'converter' function to keys of the map
// func (k StreamIter[K, V, M]) ConvertKey(by func(K) K) Iter[K, V, M] {
// 	return Stream(Convert(k.next, convert.Key[V](by)).Next, k.collector)
// }

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (k Iter[K, V, M]) FilterValue(predicate func(V) bool) Iter[K, V, M] {
	return New(loop.Filter(k.next, filter.Value[K](predicate)).Next, k.collector)
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (k Iter[K, V, M]) FiltValue(predicate func(V) (bool, error)) Iter[K, V, M] {
	return New(loop.Filt(k.next, breakFilter.Value[K](predicate)).Next, k.collector)
}

// // ConvertValue returns a stream that applies the 'converter' function to values of the map
// func (k StreamIter[K, V, M]) ConvertValue(by func(V) V) Iter[K, V, M] {
// 	return Stream(Convert(k.next, convert.Value[K](by)).Next, k.collector)
// }

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (k Iter[K, V, M]) Filter(predicate func(K, V) bool) Iter[K, V, M] {
	return New(loop.Filter(k.next, predicate).Next, k.collector)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (k Iter[K, V, M]) Filt(predicate func(K, V) (bool, error)) Iter[K, V, M] {
	return New(loop.Filt(k.next, predicate).Next, k.collector)
}

// // Convert returns a stream that applies the 'converter' function to the collection elements
// func (k StreamIter[K, V, M]) Convert(converter func(K, V) (K, V)) Iter[K, V, M] {
// 	return Stream(Convert(k.next, converter).Next, k.collector)
// }

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (k Iter[K, V, M]) Track(tracker func(K, V) error) error {
	return breakLoop.Track(k.next, tracker)
}

// Reduce reduces the key/value pairs into an one pair using the 'merge' function
func (k Iter[K, V, M]) Reduce(by func(K, V, K, V) (K, V, error)) (K, V, error) {
	return loop.Reduce(k.next, by)
}

// HasAny finds the first key/value pari that satisfies the 'predicate' function condition and returns true if successful
func (k Iter[K, V, M]) HasAny(predicate func(K, V) (bool, error)) (bool, error) {
	next := k.next
	return loop.HasAny(next, predicate)
}

// Begin creates iterator
func (k Iter[K, V, M]) Begin() c.KVIteratorBreakable[K, V] {
	return k
}

// Map collects the key/value pairs to a map
func (k Iter[K, V, M]) Map() (M, error) {
	return k.collector(k.next)
}

// MapCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type MapCollector[K comparable, V any, M map[K]V | map[K][]V] func(next func() (K, V, bool, error)) (M, error)
