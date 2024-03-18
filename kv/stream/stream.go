package stream

import (
	breakKvLoop "github.com/m4gshm/gollections/break/kv/loop"
	"github.com/m4gshm/gollections/break/kv/stream"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	breakFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/kv"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
)

// New instantiates StreamIter
func New[K comparable, V any, M map[K]V | map[K][]V](next func() (K, V, bool), collector MapCollector[K, V, M]) Iter[K, V, M] {
	return Iter[K, V, M]{next: next, collector: collector}
}

// Iter is the key/value Iterator based stream implementation.
type Iter[K comparable, V any, M map[K]V | map[K][]V] struct {
	next      func() (K, V, bool)
	collector MapCollector[K, V, M]
}

var (
	_ kv.Iterator[string, any]                                                   = (*Iter[string, any, map[string]any])(nil)
	_ Stream[string, any, Iter[string, any, map[string]any], map[string]any]     = (*Iter[string, any, map[string]any])(nil)
	_ Stream[string, any, Iter[string, any, map[string][]any], map[string][]any] = (*Iter[string, any, map[string][]any])(nil)

	_ kv.Iterator[string, any]                                                   = Iter[string, any, map[string]any]{}
	_ Stream[string, any, Iter[string, any, map[string]any], map[string]any]     = Iter[string, any, map[string]any]{}
	_ Stream[string, any, Iter[string, any, map[string][]any], map[string][]any] = Iter[string, any, map[string][]any]{}
	_ kv.IterFor[string, any, Iter[string, any, map[string][]any]]               = Iter[string, any, map[string][]any]{}
)

// Next implements kv.KVIterator
func (i Iter[K, V, M]) Next() (K, V, bool) {
	return i.next()
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FilterKey(predicate func(K) bool) Iter[K, V, M] {
	return New(kvloop.Filter(i.next, filter.Key[V](predicate)).Next, i.collector)
}

// FiltKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FiltKey(predicate func(K) (bool, error)) stream.Iter[K, V, M] {
	return stream.New(breakKvLoop.Filt(breakKvLoop.From(i.next), breakFilter.Key[V](predicate)).Next, collect(i.collector))
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (i Iter[K, V, M]) ConvertKey(by func(K) K) Iter[K, V, M] {
	return New(kvloop.Convert(i.next, convert.Key[V](by)).Next, i.collector)
}

// ConvKey returns a stream that applies the 'converter' function to keys of the map
func (i Iter[K, V, M]) ConvKey(by func(K) (K, error)) stream.Iter[K, V, M] {
	return stream.New(breakKvLoop.Conv(breakKvLoop.From(i.next), breakMapConvert.Key[V](by)).Next, collect(i.collector))
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FilterValue(predicate func(V) bool) Iter[K, V, M] {
	return New(kvloop.Filter(i.next, filter.Value[K](predicate)).Next, i.collector)
}

// FiltValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (i Iter[K, V, M]) FiltValue(predicate func(V) (bool, error)) stream.Iter[K, V, M] {
	return stream.New(breakKvLoop.Filt(breakKvLoop.From(i.next), breakFilter.Value[K](predicate)).Next, collect(i.collector))
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (i Iter[K, V, M]) ConvertValue(converter func(V) V) Iter[K, V, M] {
	return New(kvloop.Convert(i.next, convert.Value[K](converter)).Next, i.collector)
}

// ConvValue returns a stream that applies the 'converter' function to values of the map
func (i Iter[K, V, M]) ConvValue(converter func(V) (V, error)) stream.Iter[K, V, M] {
	return stream.New(breakKvLoop.Conv(breakKvLoop.From(i.next), breakMapConvert.Value[K](converter)).Next, collect(i.collector))
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (i Iter[K, V, M]) Filter(predicate func(K, V) bool) Iter[K, V, M] {
	return New(kvloop.Filter(i.next, predicate).Next, i.collector)
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (i Iter[K, V, M]) Filt(predicate func(K, V) (bool, error)) stream.Iter[K, V, M] {
	return stream.New(breakKvLoop.Filt(breakKvLoop.From(i.next), predicate).Next, collect(i.collector))
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (i Iter[K, V, M]) Convert(converter func(K, V) (K, V)) Iter[K, V, M] {
	return New(kvloop.Convert(i.next, converter).Next, i.collector)
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (i Iter[K, V, M]) Conv(converter func(K, V) (K, V, error)) stream.Iter[K, V, M] {
	return stream.New(breakKvLoop.Conv(breakKvLoop.From(i.next), converter).Next, collect(i.collector))
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (i Iter[K, V, M]) Track(tracker func(K, V) error) error {
	return loop.Track(i.next, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (i Iter[K, V, M]) TrackEach(tracker func(K, V)) {
	loop.TrackEach(i.next, tracker)
}

// Reduce reduces the key/value pairs into an one pair using the 'merge' function
func (i Iter[K, V, M]) Reduce(by func(K, K, V, V) (K, V)) (K, V) {
	return kvloop.Reduce(i.next, by)
}

// HasAny finds the first key/value pari that satisfies the 'predicate' function condition and returns true if successful
func (i Iter[K, V, M]) HasAny(predicate func(K, V) bool) bool {
	next := i.next
	return kvloop.HasAny(next, predicate)
}

// Iter creates an iterator and returns as interface
func (i Iter[K, V, M]) Iter() Iter[K, V, M] {
	return i
}

// Map collects the key/value pairs to a map
func (i Iter[K, V, M]) Map() M {
	return i.collector(i.next)
}

// Start is used with for loop construct like 'for i, k, v, ok := i.Start(); ok; k, v, ok = i.Next() { }'
func (i Iter[K, V, M]) Start() (Iter[K, V, M], K, V, bool) {
	k, v, ok := i.Next()
	return i, k, v, ok
}

// MapCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type MapCollector[K comparable, V any, M map[K]V | map[K][]V] func(next func() (K, V, bool)) M

func collect[K comparable, V any, M map[K]V | map[K][]V](collector MapCollector[K, V, M]) stream.MapCollector[K, V, M] {
	return func(next func() (K, V, bool, error)) (M, error) {
		var loopErr error
		breakKvLoop.To(next, func(err error) { loopErr = err })
		return collector(breakKvLoop.To(next, func(err error) { loopErr = err })), loopErr
	}
}
