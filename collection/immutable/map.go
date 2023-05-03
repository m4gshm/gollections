package immutable

import (
	"fmt"
	"sort"

	breakLoop "github.com/m4gshm/gollections/break/kv/loop"
	breakKvStream "github.com/m4gshm/gollections/break/kv/stream"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	breakMapFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/kv"
	"github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/slice"
)

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) Map[K, V] {
	return Map[K, V]{elements: elements}
}

// Map is a collection that provides elements access by an unique key.
type Map[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ collection.Map[int, any]                         = (*Map[int, any])(nil)
	_ collection.Map[int, any]                         = Map[int, any]{}
	_ loop.Looper[int, any, *map_.Iter[int, any]]      = (*Map[int, any])(nil)
	_ loop.Looper[int, any, *map_.Iter[int, any]]      = Map[int, any]{}
	_ c.KeyVal[MapKeys[int, any], MapValues[int, any]] = (*Map[int, any])(nil)
	_ c.KeyVal[MapKeys[int, any], MapValues[int, any]] = Map[int, any]{}
	_ fmt.Stringer                                     = (*Map[int, any])(nil)
	_ fmt.Stringer                                     = Map[int, any]{}
)

// Iter creates an iterator and returns as interface
func (m Map[K, V]) Iter() kv.Iterator[K, V] {
	h := m.Head()
	return &h
}

// Loop creates an iterator and returns as implementation type reference
func (m Map[K, V]) Loop() *map_.Iter[K, V] {
	h := m.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (m Map[K, V]) Head() map_.Iter[K, V] {
	return map_.NewIter(m.elements)
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m Map[K, V]) First() (map_.Iter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Map collects the key/value pairs to a map
func (m Map[K, V]) Map() map[K]V {
	return map_.Clone(m.elements)
}

// Len returns amount of elements
func (m Map[K, V]) Len() int {
	return len(m.elements)
}

// IsEmpty returns true if the map is empty
func (m Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Contains checks is the map contains a key
func (m Map[K, V]) Contains(key K) (ok bool) {
	if m.elements != nil {
		_, ok = m.elements[key]
	}
	return ok
}

// Get returns the value for a key.
// If ok==false, then the map does not contain the key.
func (m Map[K, V]) Get(key K) (element V, ok bool) {
	if m.elements != nil {
		element, ok = m.elements[key]
	}
	return element, ok
}

// Keys resutrns keys collection
func (m Map[K, V]) Keys() MapKeys[K, V] {
	return m.K()
}

// K resutrns keys collection impl
func (m Map[K, V]) K() MapKeys[K, V] {
	return WrapKeys(m.elements)
}

// Values resutrns values collection
func (m Map[K, V]) Values() MapValues[K, V] {
	return m.V()
}

// V resutrns values collection impl
func (m Map[K, V]) V() MapValues[K, V] {
	return WrapVal(m.elements)
}

// Sort sorts the keys
func (m Map[K, V]) Sort(less slice.Less[K]) ordered.Map[K, V] {
	return m.sortBy(sort.Slice, less)
}

// StableSort sorts the keys
func (m Map[K, V]) StableSort(less slice.Less[K]) ordered.Map[K, V] {
	return m.sortBy(sort.SliceStable, less)
}

func (m Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) ordered.Map[K, V] {
	return ordered.WrapMap(slice.Sort(map_.Keys(m.elements), sorter, less), m.elements)
}

func (m Map[K, V]) String() string {
	return map_.ToString(m.elements)
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (m Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.Track(m.elements, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (m Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEach(m.elements, tracker)
}

// For applies the 'walker' function for every key/value pair. Return the c.ErrBreak to stop.
func (m Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return map_.For(m.elements, walker)
}

// ForEach applies the 'walker' function for every key/value pair
func (m Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	map_.ForEach(m.elements, walker)
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterKey(predicate func(K) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Key[V](predicate)).Next, loop.ToMap[K, V])
}

// FiltKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Key[V](predicate)).Next, breakLoop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(by func(K) K) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Key[V](by)).Next, loop.ToMap[K, V])
}

// ConvKey returns a stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Key[V](converter)).Next, breakLoop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Value[K](predicate)).Next, loop.ToMap[K, V])
}

// FiltValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltValue(predicate func(V) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Value[K](predicate)).Next, breakLoop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(by func(V) V) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Value[K](by)).Next, loop.ToMap[K, V])
}

// ConvValue returns a stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Value[K](converter)).Next, breakLoop.ToMap[K, V])
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, predicate).Next, loop.ToMap[K, V])
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filt(predicate func(K, V) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next, breakLoop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, converter).Next, loop.ToMap[K, V])
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (m Map[K, V]) Conv(converter func(K, V) (K, V, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), converter).Next, breakLoop.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m Map[K, V]) Reduce(merge func(K, V, K, V) (K, V)) (K, V) {
	return map_.Reduce(m.elements, merge)
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	return map_.HasAny(m.elements, predicate)
}
