package ordered

import (
	"fmt"

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

// WrapMap instantiates an ordered Map using a map and an order slice as internal storage
func WrapMap[K comparable, V any](order []K, elements map[K]V) *Map[K, V] {
	return &Map[K, V]{order: order, elements: elements}
}

// Map is a collection implementation that provides elements access by an unique key.
type Map[K comparable, V any] struct {
	order    []K
	elements map[K]V
}

var (
	_ c.Settable[int, any]                                        = (*Map[int, any])(nil)
	_ c.SettableNew[int, any]                                     = (*Map[int, any])(nil)
	_ c.SettableMap[c.TrackEachLoop[int, any]]                    = (*Map[int, any])(nil)
	_ c.ImmutableMapConvert[ordered.Map[int, any]]                = (*Map[int, any])(nil)
	_ collection.Map[int, any]                                    = (*Map[int, any])(nil)
	_ loop.Looper[int, any, *ordered.MapIter[int, any]]           = (*Map[int, any])(nil)
	_ c.KeyVal[ordered.MapKeys[int], ordered.MapValues[int, any]] = (*Map[int, any])(nil)
	_ fmt.Stringer                                                = (*Map[int, any])(nil)
)

// Iter creates an iterator and returns as interface
func (m *Map[K, V]) Iter() kv.Iterator[K, V] {
	h := m.Head()
	return &h
}

// Loop creates an iterator and returns as implementation type reference
func (m *Map[K, V]) Loop() *ordered.MapIter[K, V] {
	h := m.Head()
	return &h
}

// Head creates an iterator and returns as implementation type value
func (m *Map[K, V]) Head() ordered.MapIter[K, V] {
	var (
		order    []K
		elements map[K]V
	)
	if m != nil {
		elements = m.elements
		order = m.order
	}
	return ordered.NewMapIter(elements, slice.NewHead(order))
}

// Tail creates an iterator pointing to the end of the collection
func (m *Map[K, V]) Tail() ordered.MapIter[K, V] {
	var (
		order    []K
		elements map[K]V
	)
	if m != nil {
		elements = m.elements
		order = m.order
	}
	return ordered.NewMapIter(elements, slice.NewTail(order))
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m *Map[K, V]) First() (ordered.MapIter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Map collects the key/value pairs to a map
func (m *Map[K, V]) Map() map[K]V {
	if m == nil {
		return nil
	}
	return map_.Clone(m.elements)
}

// Sort sorts keys in-place (no copy)
func (m *Map[K, V]) Sort(comparer slice.Comparer[K]) *Map[K, V] {
	return m.sortBy(slice.Sort, comparer)
}

// StableSort sorts keys in-place (no copy)
func (m *Map[K, V]) StableSort(comparer slice.Comparer[K]) *Map[K, V] {
	return m.sortBy(slice.StableSort, comparer)
}

func (m *Map[K, V]) sortBy(sorter func([]K, slice.Comparer[K]) []K, comparer slice.Comparer[K]) *Map[K, V] {
	if m != nil {
		sorter(m.order, comparer)
	}
	return m
}

// Len returns the amount of elements contained in the map
func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.order)
}

// IsEmpty returns true if the map is empty
func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// For applies the 'walker' function for key/value pairs. Return the c.ErrBreak to stop.
func (m *Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	if m == nil {
		return nil
	}
	return map_.ForOrdered(m.order, m.elements, walker)
}

// ForEach applies the 'walker' function for every key/value pair
func (m *Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	if m == nil {
		return
	}
	map_.ForEachOrdered(m.order, m.elements, walker)
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (m *Map[K, V]) Track(tracker func(K, V) error) error {
	if m == nil {
		return nil
	}
	return map_.TrackOrdered(m.order, m.elements, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (m *Map[K, V]) TrackEach(tracker func(K, V)) {
	if m == nil {
		return
	}
	map_.TrackEachOrdered(m.order, m.elements, tracker)
}

// Contains checks is the map contains a key
func (m *Map[K, V]) Contains(key K) bool {
	if m == nil {
		return false
	}
	_, ok := m.elements[key]
	return ok
}

// Get returns the value for a key.
// If ok==false, then the map does not contain the key.
func (m *Map[K, V]) Get(key K) (V, bool) {
	if m == nil {
		var z V
		return z, false
	}
	val, ok := m.elements[key]
	return val, ok
}

// Set sets the value for a key
func (m *Map[K, V]) Set(key K, value V) {
	if m == nil {
		return
	}
	u := m.elements
	if u == nil {
		u = map[K]V{}
		m.elements = u
	}
	if _, ok := u[key]; !ok {
		m.order = append(m.order, key)
	}
	u[key] = value
}

// SetNew sets the value fo a key only if the key is not exists in the map
func (m *Map[K, V]) SetNew(key K, value V) bool {
	if m == nil {
		return false
	}
	u := m.elements
	if u == nil {
		u = map[K]V{}
		m.elements = u
	}

	if _, ok := u[key]; !ok {
		u[key] = value
		m.order = append(m.order, key)
		return true
	}
	return false
}

// Keys resutrns keys collection
func (m *Map[K, V]) Keys() ordered.MapKeys[K] {
	var order []K
	if m != nil {
		order = m.order
	}
	return ordered.WrapKeys(order)
}

// Values resutrns values collection
func (m *Map[K, V]) Values() ordered.MapValues[K, V] {
	var (
		order    []K
		elements map[K]V
	)
	if m != nil {
		order, elements = m.order, m.elements
	}
	return ordered.WrapVal(order, elements)
}

// String string representation on the map
func (m *Map[K, V]) String() string {
	var (
		order    []K
		elements map[K]V
	)
	if m != nil {
		order, elements = m.order, m.elements
	}
	return map_.ToStringOrdered(order, elements)
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterKey(predicate func(K) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Key[V](predicate)).Next, loop.ToMap[K, V])
}

// FiltKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Key[V](predicate)).Next, breakLoop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m *Map[K, V]) ConvertKey(converter func(K) K) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Key[V](converter)).Next, loop.ToMap[K, V])
}

// ConvKey returns a stream that applies the 'converter' function to keys of the map
func (m *Map[K, V]) ConvKey(converter func(K) (K, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Key[V](converter)).Next, breakLoop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterValue(predicate func(V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Value[K](predicate)).Next, loop.ToMap[K, V])
}

// FiltValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FiltValue(predicate func(V) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Value[K](predicate)).Next, breakLoop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m *Map[K, V]) ConvertValue(converter func(V) V) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Value[K](converter)).Next, loop.ToMap[K, V])
}

// ConvValue returns a stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Value[K](converter)).Next, breakLoop.ToMap[K, V])
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filter(predicate func(K, V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, predicate).Next, loop.ToMap[K, V])
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filt(predicate func(K, V) (bool, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate).Next, breakLoop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Convert(converter func(K, V) (K, V)) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, converter).Next, loop.ToMap[K, V])
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Conv(converter func(K, V) (K, V, error)) breakKvStream.Iter[K, V, map[K]V] {
	h := m.Head()
	return breakKvStream.New(breakLoop.Conv(breakLoop.From(h.Next), converter).Next, breakLoop.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m *Map[K, V]) Reduce(merge func(K, K, V, V) (K, V)) (k K, v V) {
	if m != nil {
		k, v = map_.Reduce(m.elements, merge)
	}
	return k, v
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m *Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	return map_.HasAny(m.elements, predicate)
}

// Immutable converts to an immutable map instance
func (m *Map[K, V]) Immutable() ordered.Map[K, V] {
	var e map[K]V
	var o []K
	if m != nil {
		e = map_.Clone(m.elements)
		o = slice.Clone(m.order)
	}
	return ordered.WrapMap(o, e)
}

// SetMap inserts all elements from the 'other' map
func (m *Map[K, V]) SetMap(kvs c.TrackEachLoop[K, V]) {
	if m == nil || kvs == nil {
		return
	}
	kvs.TrackEach(func(key K, value V) { m.Set(key, value) })
}

func addToMap[K comparable, V any](key K, val V, order []K, uniques map[K]V) []K {
	if _, ok := uniques[key]; !ok {
		order = append(order, key)
		uniques[key] = val
	}
	return order
}
