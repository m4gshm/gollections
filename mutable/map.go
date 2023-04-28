package mutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/kv/loop"
	breakLoop "github.com/m4gshm/gollections/kv/loop/break/loop"
	"github.com/m4gshm/gollections/map_"
	breakMapConvert "github.com/m4gshm/gollections/map_/break/map_/convert"
	breakMapFilter "github.com/m4gshm/gollections/map_/break/map_/filter"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/map_/iter"
)

// NewMap instantiates Map with a predefined capacity.
func NewMap[K comparable, V any](capacity int) *Map[K, V] {
	return WrapMap(make(map[K]V, capacity))
}

// NewMapKV converts a slice of key/value pairs into a Map instance.
func NewMapKV[K comparable, V any](elements []c.KV[K, V]) *Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// ToMap instantiates Map and copies elements to it.
func ToMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	m := Map[K, V](elements)
	return &m
}

// Map is the Collection implementation based on the embedded map.
type Map[K comparable, V any] map[K]V

var (
	_ c.Deleteable[int]                                        = (*Map[int, any])(nil)
	_ c.Removable[int, any]                                    = (*Map[int, any])(nil)
	_ c.Settable[int, any]                                     = (*Map[int, any])(nil)
	_ c.SettableNew[int, any]                                  = (*Map[int, any])(nil)
	_ c.SettableMap[int, any]                                  = (*Map[int, any])(nil)
	_ c.ImmutableMapConvert[int, any, immutable.Map[int, any]] = (*Map[int, any])(nil)
	_ c.Map[int, any]                                          = (*Map[int, any])(nil)
	_ fmt.Stringer                                             = (*Map[int, any])(nil)
)

// Begin creates iterator
func (m *Map[K, V]) Begin() c.KVIterator[K, V] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m *Map[K, V]) Head() iter.MapIter[K, V] {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return iter.New(out)
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m *Map[K, V]) First() (iter.MapIter[K, V], K, V, bool) {
	var out map[K]V
	if m != nil {
		out = *m
	}
	var (
		iterator           = iter.New(out)
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Map collects the key/value pairs to a map
func (m *Map[K, V]) Map() (out map[K]V) {
	if m == nil {
		return
	}
	return map_.Clone(*m)
}

// Len returns amount of elements
func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(*m)
}

// IsEmpty returns true if the map is empty
func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// For applies the 'walker' function for every key/value pair. Return the c.ErrBreak to stop.
func (m *Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	if m == nil {
		return nil
	}
	return map_.For(*m, walker)
}

// ForEach applies the 'walker' function for every element
func (m *Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	if m != nil {
		map_.ForEach(*m, walker)
	}
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (m *Map[K, V]) Track(tracker func(K, V) error) error {
	if m == nil {
		return nil
	}
	return map_.Track(*m, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (m *Map[K, V]) TrackEach(tracker func(K, V)) {
	if m != nil {
		map_.TrackEach(*m, tracker)
	}
}

// Contains checks is the map contains a key
func (m *Map[K, V]) Contains(key K) (ok bool) {
	if m != nil {
		_, ok = (*m)[key]
	}
	return ok
}

// Get returns the value for a key.
// If ok==false, then the map does not contain the key.
func (m *Map[K, V]) Get(key K) (val V, ok bool) {
	if m != nil {
		val, ok = (*m)[key]
	}
	return val, ok
}

// Set sets the value for a key
func (m *Map[K, V]) Set(key K, value V) {
	if m == nil {
		return
	} else if (*m) == nil {
		*m = Map[K, V]{}
	}
	(*m)[key] = value
}

// SetNew sets the value fo a key only if the key is not exists in the map
func (m *Map[K, V]) SetNew(key K, value V) bool {
	if m == nil {
		return false
	} else if (*m) == nil {
		*m = Map[K, V]{}
	}

	if _, ok := (*m)[key]; !ok {
		(*m)[key] = value
		return true
	}
	return false
}

// Delete removes value by their keys from the map
func (m *Map[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		m.DeleteOne(key)
	}
}

// DeleteOne removes an value by the key from the map
func (m *Map[K, V]) DeleteOne(key K) {
	if m != nil {
		delete(*m, key)
	}
}

// Remove removes value by key and return it
func (m *Map[K, V]) Remove(key K) (v V, ok bool) {
	if m == nil {
		return v, ok
	}
	v, ok = m.Get(key)
	m.Delete(key)
	return v, ok
}

// Keys resutrns keys collection
func (m *Map[K, V]) Keys() c.Collection[K] {
	return m.K()
}

// K resutrns keys collection impl
func (m *Map[K, V]) K() immutable.MapKeys[K, V] {
	var elements map[K]V
	if m != nil {
		elements = *m
	}
	return immutable.WrapKeys(elements)
}

// Values resutrns values collection
func (m *Map[K, V]) Values() c.Collection[V] {
	return m.V()
}

// V resutrns values collection impl
func (m *Map[K, V]) V() immutable.MapValues[K, V] {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return immutable.WrapVal(out)
}

// String string representation on the map
func (m *Map[K, V]) String() string {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return map_.ToString(out)
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterKey(predicate func(K) bool) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Filter(h.Next, filter.Key[V](predicate)).Next, loop.ToMap[K, V])
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) c.KVStreamBreakable[K, V, map[K]V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Key[V](predicate)).Next, breakLoop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m *Map[K, V]) ConvertKey(converter func(K) K) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Convert(h.Next, convert.Key[V](converter)).Next, loop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) c.KVStreamBreakable[K, V, map[K]V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Key[V](converter)).Next, breakLoop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterValue(predicate func(V) bool) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Filter(h.Next, filter.Value[K](predicate)).Next, loop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FiltValue(predicate func(V) (bool, error)) c.KVStreamBreakable[K, V, map[K]V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Value[K](predicate)).Next, breakLoop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m *Map[K, V]) ConvertValue(converter func(V) V) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Convert(h.Next, convert.Value[K](converter)).Next, loop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) c.KVStreamBreakable[K, V, map[K]V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Value[K](converter)).Next, breakLoop.ToMap[K, V])
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filter(predicate func(K, V) bool) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Filter(h.Next, predicate).Next, loop.ToMap[K, V])
}

// Filt returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filt(predicate func(K, V) (bool, error)) c.KVStreamBreakable[K, V, map[K]V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next, breakLoop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Convert(converter func(K, V) (K, V)) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Convert(h.Next, converter).Next, loop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Conv(converter func(K, V) (K, V, error)) c.KVStreamBreakable[K, V, map[K]V] {
	h := m.Head()
	return breakLoop.Stream(breakLoop.Conv(breakLoop.From(h.Next), converter).Next, breakLoop.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m *Map[K, V]) Reduce(merge func(K, V, K, V) (K, V)) (k K, v V) {
	if m != nil {
		k, v = map_.Reduce(*m, merge)
	}
	return k, v
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m *Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	if m != nil {
		return map_.HasAny(*m, predicate)
	}
	return false
}

// Immutable converts to an immutable map instance
func (m *Map[K, V]) Immutable() immutable.Map[K, V] {
	return immutable.WrapMap(m.Map())
}

// SetMap inserts all elements from the 'other' map
func (m *Map[K, V]) SetMap(other c.Map[K, V]) {
	if m == nil || other == nil {
		return
	}
	other.TrackEach(func(key K, value V) { m.Set(key, value) })
}
