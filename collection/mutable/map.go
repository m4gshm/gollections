package mutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/kv/loop"
	breakKvStream "github.com/m4gshm/gollections/break/kv/stream"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	breakMapFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable"
	"github.com/m4gshm/gollections/collection/mutable/ordered"
	"github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/slice"
)

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	m := Map[K, V](elements)
	return &m
}

// Map is a collection implementation that provides elements access by an unique key.
type Map[K comparable, V any] map[K]V

var (
	_ c.Deleteable[int]                                                    = (*Map[int, any])(nil)
	_ c.Removable[int, any]                                                = (*Map[int, any])(nil)
	_ c.Settable[int, any]                                                 = (*Map[int, any])(nil)
	_ c.SettableNew[int, any]                                              = (*Map[int, any])(nil)
	_ c.SettableMap[c.TrackEach[int, any]]                                 = (*Map[int, any])(nil)
	_ c.ImmutableMapConvert[immutable.Map[int, any]]                       = (*Map[int, any])(nil)
	_ collection.Map[int, any]                                             = (*Map[int, any])(nil)
	_ c.KeyVal[immutable.MapKeys[int, any], immutable.MapValues[int, any]] = (*Map[int, any])(nil)
	_ fmt.Stringer                                                         = (*Map[int, any])(nil)
)

// Iter creates an iterator and returns as interface
func (m *Map[K, V]) Loop() loop.Loop[K, V] {
	h := m.Head()
	return h.Next
}

// Head creates an iterator and returns as implementation type value
func (m *Map[K, V]) Head() map_.Iter[K, V] {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return map_.NewIter(out)
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m *Map[K, V]) First() (map_.Iter[K, V], K, V, bool) {
	var out map[K]V
	if m != nil {
		out = *m
	}
	var (
		iterator           = map_.NewIter(out)
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

// Sort sorts keys in-place (no copy)
func (m *Map[K, V]) Sort(comparer slice.Comparer[K]) *ordered.Map[K, V] {
	return m.sortBy(slice.Sort, comparer)
}

// StableSort sorts keys in-place (no copy)
func (m *Map[K, V]) StableSort(comparer slice.Comparer[K]) *ordered.Map[K, V] {
	return m.sortBy(slice.StableSort, comparer)
}

func (m *Map[K, V]) sortBy(sorter func([]K, slice.Comparer[K]) []K, comparer slice.Comparer[K]) *ordered.Map[K, V] {
	if m != nil {
		return ordered.NewMapOf(sorter(m.Keys().Slice(), comparer), m.Map())
	}
	return nil
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
func (m *Map[K, V]) Keys() immutable.MapKeys[K, V] {
	var elements map[K]V
	if m != nil {
		elements = *m
	}
	return immutable.WrapKeys(elements)
}

// Values resutrns values collection
func (m *Map[K, V]) Values() immutable.MapValues[K, V] {
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
func (m *Map[K, V]) FilterKey(predicate func(K) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Key[V](predicate)), loop.ToMap[K, V])
}

// FiltKey returns a breakable stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Key[V](predicate)), breakLoop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m *Map[K, V]) ConvertKey(converter func(K) K) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Key[V](converter)), loop.ToMap[K, V])
}

// ConvKey returns a breabkable stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Key[V](converter)), breakLoop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterValue(predicate func(V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Value[K](predicate)), loop.ToMap[K, V])
}

// FiltValue returns a breakable stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FiltValue(predicate func(V) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Value[K](predicate)), breakLoop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m *Map[K, V]) ConvertValue(converter func(V) V) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Value[K](converter)), loop.ToMap[K, V])
}

// ConvValue returns a breakable stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Value[K](converter)), breakLoop.ToMap[K, V])
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filter(predicate func(K, V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, predicate), loop.ToMap[K, V])
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filt(predicate func(K, V) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate), breakLoop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Convert(converter func(K, V) (K, V)) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, converter), loop.ToMap[K, V])
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Conv(converter func(K, V) (K, V, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), converter), breakLoop.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m *Map[K, V]) Reduce(merge func(K, K, V, V) (K, V)) (k K, v V) {
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
func (m *Map[K, V]) SetMap(other c.TrackEach[K, V]) {
	if m == nil || other == nil {
		return
	}
	other.TrackEach(func(key K, value V) { m.Set(key, value) })
}
