package loop

import (
	breakKvLoop "github.com/m4gshm/gollections/break/kv/loop"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	breakFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
)

// Stream instantiates Iterator wrapper that converts the elements into key/value pairs and iterates over them.
func Stream[K comparable, V any, M map[K]V | map[K][]V](next func() (K, V, bool), collector MapCollector[K, V, M]) StreamIter[K, V, M] {
	return StreamIter[K, V, M]{next: next, collector: collector}
}

// StreamIter is the key/value Iterator based stream implementation.
type StreamIter[K comparable, V any, M map[K]V | map[K][]V] struct {
	next      func() (K, V, bool)
	collector MapCollector[K, V, M]
}

var (
	_ c.KVIterator[string, any]                 = (*StreamIter[string, any, map[string]any])(nil)
	_ c.KVStream[string, any, map[string]any]   = (*StreamIter[string, any, map[string]any])(nil)
	_ c.KVStream[string, any, map[string][]any] = (*StreamIter[string, any, map[string][]any])(nil)

	_ c.KVIterator[string, any]                 = StreamIter[string, any, map[string]any]{}
	_ c.KVStream[string, any, map[string]any]   = StreamIter[string, any, map[string]any]{}
	_ c.KVStream[string, any, map[string][]any] = StreamIter[string, any, map[string][]any]{}
)

// Next implements c.KVIterator
func (k StreamIter[K, V, M]) Next() (K, V, bool) {
	return k.next()
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (k StreamIter[K, V, M]) FilterKey(predicate func(K) bool) c.KVStream[K, V, M] {
	return Stream(Filter(k.next, filter.Key[V](predicate)).Next, k.collector)
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (k StreamIter[K, V, M]) FiltKey(predicate func(K) (bool, error)) c.KVStreamBreakable[K, V, M] {
	return breakKvLoop.Stream(breakKvLoop.Filt(breakKvLoop.From(k.next), breakFilter.Key[V](predicate)).Next, collect(k.collector))
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (k StreamIter[K, V, M]) ConvertKey(by func(K) K) c.KVStream[K, V, M] {
	return Stream(Convert(k.next, convert.Key[V](by)).Next, k.collector)
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (k StreamIter[K, V, M]) ConvKey(by func(K) (K, error)) c.KVStreamBreakable[K, V, M] {
	return breakKvLoop.Stream(breakKvLoop.Conv(breakKvLoop.From(k.next), breakMapConvert.Key[V](by)).Next, collect(k.collector))
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (k StreamIter[K, V, M]) FilterValue(predicate func(V) bool) c.KVStream[K, V, M] {
	return Stream(Filter(k.next, filter.Value[K](predicate)).Next, k.collector)
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (k StreamIter[K, V, M]) FiltValue(predicate func(V) (bool, error)) c.KVStreamBreakable[K, V, M] {
	return breakKvLoop.Stream(breakKvLoop.Filt(breakKvLoop.From(k.next), breakFilter.Value[K](predicate)).Next, collect(k.collector))
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (k StreamIter[K, V, M]) ConvertValue(converter func(V) V) c.KVStream[K, V, M] {
	return Stream(Convert(k.next, convert.Value[K](converter)).Next, k.collector)
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (k StreamIter[K, V, M]) ConvValue(converter func(V) (V, error)) c.KVStreamBreakable[K, V, M] {
	return breakKvLoop.Stream(breakKvLoop.Conv(breakKvLoop.From(k.next), breakMapConvert.Value[K](converter)).Next, collect(k.collector))
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (k StreamIter[K, V, M]) Filter(predicate func(K, V) bool) c.KVStream[K, V, M] {
	return Stream(Filter(k.next, predicate).Next, k.collector)
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (k StreamIter[K, V, M]) Filt(predicate func(K, V) (bool, error)) c.KVStreamBreakable[K, V, M] {
	return breakKvLoop.Stream(breakKvLoop.Filt(breakKvLoop.From(k.next), predicate).Next, collect(k.collector))
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (k StreamIter[K, V, M]) Convert(converter func(K, V) (K, V)) c.KVStream[K, V, M] {
	return Stream(Convert(k.next, converter).Next, k.collector)
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (k StreamIter[K, V, M]) Conv(converter func(K, V) (K, V, error)) c.KVStreamBreakable[K, V, M] {
	return breakKvLoop.Stream(breakKvLoop.Conv(breakKvLoop.From(k.next), converter).Next, collect(k.collector))
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (k StreamIter[K, V, M]) Track(tracker func(K, V) error) error {
	return loop.Track(k.next, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (k StreamIter[K, V, M]) TrackEach(tracker func(K, V)) {
	loop.TrackEach(k.next, tracker)
}

// Reduce reduces the key/value pairs into an one pair using the 'merge' function
func (k StreamIter[K, V, M]) Reduce(by func(K, V, K, V) (K, V)) (K, V) {
	return Reduce(k.next, by)
}

// HasAny finds the first key/value pari that satisfies the 'predicate' function condition and returns true if successful
func (k StreamIter[K, V, M]) HasAny(predicate func(K, V) bool) bool {
	next := k.next
	return HasAny(next, predicate)
}

// Begin creates iterator
func (k StreamIter[K, V, M]) Begin() c.KVIterator[K, V] {
	return k
}

// Map collects the key/value pairs to a map
func (k StreamIter[K, V, M]) Map() M {
	return k.collector(k.next)
}

// MapCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type MapCollector[K comparable, V any, M map[K]V | map[K][]V] func(next func() (K, V, bool)) M

func collect[K comparable, V any, M map[K]V | map[K][]V](collector MapCollector[K, V, M]) breakKvLoop.MapCollector[K, V, M] {
	return func(next func() (K, V, bool, error)) (M, error) {
		var loopErr error
		breakKvLoop.To(next, func(err error) { loopErr = err })
		return collector(breakKvLoop.To(next, func(err error) { loopErr = err })), loopErr
	}
}
