package iter

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
)

// NewKVPipe instantiates Iterator wrapper that converts the elements into key/value pairs and iterates over them.
func NewKVPipe[K comparable, V any, M map[K]V | map[K][]V](next func() (K, V, bool), collector MapCollector[K, V, M]) KVIterPipe[K, V, M] {
	return KVIterPipe[K, V, M]{next: next, collector: collector}
}

// KVIterPipe is the key/value Iterator based pipe implementation.
type KVIterPipe[K comparable, V any, M map[K]V | map[K][]V] struct {
	next      func() (K, V, bool)
	collector MapCollector[K, V, M]
}

var (
	_ c.KVIterator[string, any]                     = (*KVIterPipe[string, any, map[string]any])(nil)
	_ c.MapTransform[string, any, map[string]any]   = (*KVIterPipe[string, any, map[string]any])(nil)
	_ c.MapTransform[string, any, map[string][]any] = (*KVIterPipe[string, any, map[string][]any])(nil)

	_ c.KVIterator[string, any]                     = KVIterPipe[string, any, map[string]any]{}
	_ c.MapTransform[string, any, map[string]any]   = KVIterPipe[string, any, map[string]any]{}
	_ c.MapTransform[string, any, map[string][]any] = KVIterPipe[string, any, map[string][]any]{}
)

// Next implements c.KVIterator
func (k KVIterPipe[K, V, M]) Next() (K, V, bool) {
	return k.next()
}

// FilterKey returns a pipe consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (k KVIterPipe[K, V, M]) FilterKey(predicate func(K) bool) c.MapTransform[K, V, M] {
	return NewKVPipe(FilterKV(k.next, filter.Key[V](predicate)).Next, k.collector)
}

// ConvertKey returns a pipe that applies the 'converter' function to keys of the map
func (k KVIterPipe[K, V, M]) ConvertKey(by func(K) K) c.MapTransform[K, V, M] {
	return NewKVPipe(ConvertKV(k.next, convert.Key[V](by)).Next, k.collector)
}

// FilterValue returns a pipe consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (k KVIterPipe[K, V, M]) FilterValue(predicate func(V) bool) c.MapTransform[K, V, M] {
	return NewKVPipe(FilterKV(k.next, filter.Value[K](predicate)).Next, k.collector)
}

// ConvertValue returns a pipe that applies the 'converter' function to values of the map
func (k KVIterPipe[K, V, M]) ConvertValue(by func(V) V) c.MapTransform[K, V, M] {
	return NewKVPipe(ConvertKV(k.next, convert.Value[K](by)).Next, k.collector)
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (k KVIterPipe[K, V, M]) Filter(predicate func(K, V) bool) c.MapTransform[K, V, M] {
	return NewKVPipe(FilterKV(k.next, predicate).Next, k.collector)
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (k KVIterPipe[K, V, M]) Convert(converter func(K, V) (K, V)) c.MapTransform[K, V, M] {
	return NewKVPipe(ConvertKV(k.next, converter).Next, k.collector)
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (k KVIterPipe[K, V, M]) Track(tracker func(K, V) error) error {
	return loop.Track(k.next, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (k KVIterPipe[K, V, M]) TrackEach(tracker func(K, V)) {
	loop.TrackEach(k.next, tracker)
}

// Reduce reduces the key/value pairs into an one pair using the 'merge' function
func (k KVIterPipe[K, V, M]) Reduce(by func(K, V, K, V) (K, V)) (K, V) {
	return loop.ReduceKV(k.next, by)
}

// HasAny finds the first key/value pari that satisfies the 'predicate' function condition and returns true if successful
func (k KVIterPipe[K, V, M]) HasAny(predicate func(K, V) bool) bool {
	next := k.next
	return loop.HasAnyKV(next, predicate)
}

// Begin creates iterator
func (k KVIterPipe[K, V, M]) Begin() c.KVIterator[K, V] {
	return k
}

// Map collects the key/value pairs to a map
func (k KVIterPipe[K, V, M]) Map() M {
	return k.collector(k)
}

// MapCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type MapCollector[K comparable, V any, M map[K]V | map[K][]V] func(c.KVIterator[K, V]) M
