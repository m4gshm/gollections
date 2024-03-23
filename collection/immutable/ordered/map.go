package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/kv/loop"
	kvstream "github.com/m4gshm/gollections/break/kv/stream"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	breakMapFilter "github.com/m4gshm/gollections/break/map_/filter"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/kv/stream"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/slice"
)

// WrapMap instantiates ordered Map using a map and an order slice as internal storage.
func WrapMap[K comparable, V any](order []K, elements map[K]V) Map[K, V] {
	return Map[K, V]{order: order, elements: elements}
}

// Map is a collection implementation that provides elements access by an unique key.
type Map[K comparable, V any] struct {
	order    []K
	elements map[K]V
}

var (
	_ collection.Map[int, any]                    = (*Map[int, any])(nil)
	_ collection.Map[int, any]                    = Map[int, any]{}
	_ c.KeyVal[MapKeys[int], MapValues[int, any]] = (*Map[int, any])(nil)
	_ c.KeyVal[MapKeys[int], MapValues[int, any]] = Map[int, any]{}
	_ fmt.Stringer                                = (*Map[int, any])(nil)
	_ fmt.Stringer                                = Map[int, any]{}
)

// Iter creates an iterator and returns as interface
func (m Map[K, V]) Loop() loop.Loop[K, V] {
	h := m.Head()
	return h.Next
}

// Head creates an iterator and returns as implementation type value
func (m Map[K, V]) Head() MapIter[K, V] {
	return NewMapIter(m.elements, slice.NewHead(m.order))
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m Map[K, V]) First() (MapIter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Tail creates an iterator pointing to the end of the map
func (m Map[K, V]) Tail() MapIter[K, V] {
	return NewMapIter(m.elements, slice.NewTail(m.order))
}

// Map collects the key/value pairs to a map
func (m Map[K, V]) Map() map[K]V {
	return map_.Clone(m.elements)
}

// Len returns amount of elements
func (m Map[K, V]) Len() int {
	return len(m.order)
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
func (m Map[K, V]) Keys() MapKeys[K] {
	return WrapKeys(m.order)
}

// Values resutrns values collection
func (m Map[K, V]) Values() MapValues[K, V] {
	return WrapVal(m.order, m.elements)
}

// Sort returns sorted by keys map
func (m Map[K, V]) Sort(comparer slice.Comparer[K]) Map[K, V] {
	return m.sortBy(slice.Sort, comparer)
}

// StableSort returns sorted by keys map
func (m Map[K, V]) StableSort(comparer slice.Comparer[K]) Map[K, V] {
	return m.sortBy(slice.StableSort, comparer)
}

func (m Map[K, V]) sortBy(sorter func([]K, slice.Comparer[K]) []K, comparer slice.Comparer[K]) Map[K, V] {
	return WrapMap(sorter(slice.Clone(m.order), comparer), m.elements)
}

func (m Map[K, V]) String() string {
	return map_.ToStringOrdered(m.order, m.elements)
}

// Track applies the 'tracker' function for key/value pairs. Return the c.ErrBreak to stop.
func (m Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.TrackOrdered(m.order, m.elements, tracker)
}

// TrackEach applies the 'tracker' function for every key/value pairs
func (m Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEachOrdered(m.order, m.elements, tracker)
}

// For applies the 'walker' function for key/value pairs. Return the c.ErrBreak to stop.
func (m Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return map_.ForOrdered(m.order, m.elements, walker)
}

// ForEach applies the 'walker' function for every key/value pair
func (m Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	map_.ForEachOrdered(m.order, m.elements, walker)
}

// FilterKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterKey(predicate func(K) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Key[V](predicate)), loop.ToMap[K, V])
}

// FiltKey returns a stream consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) kvstream.Iter[K, V, map[K]V] {
	h := m.Head()
	return kvstream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Key[V](predicate)), breakLoop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(by func(K) K) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Key[V](by)), loop.ToMap[K, V])
}

// ConvKey returns a stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) kvstream.Iter[K, V, map[K]V] {
	h := m.Head()
	return kvstream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Key[V](converter)), breakLoop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, filter.Value[K](predicate)), loop.ToMap[K, V])
}

// FiltValue returns a breakable stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltValue(predicate func(V) (bool, error)) kvstream.Iter[K, V, map[K]V] {
	h := m.Head()
	return kvstream.New(breakLoop.Filt(breakLoop.From(h.Next), breakMapFilter.Value[K](predicate)), breakLoop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(converter func(V) V) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, convert.Value[K](converter)), loop.ToMap[K, V])
}

// ConvValue returns a breakable stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) kvstream.Iter[K, V, map[K]V] {
	h := m.Head()
	return kvstream.New(breakLoop.Conv(breakLoop.From(h.Next), breakMapConvert.Value[K](converter)), breakLoop.ToMap[K, V])
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Filter(h.Next, predicate), loop.ToMap[K, V])
}

// Filt returns a breakable stream consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filt(predicate func(K, V) (bool, error)) kvstream.Iter[K, V, map[K]V] {
	h := m.Head()
	return kvstream.New(breakLoop.Filt(breakLoop.From(h.Next), predicate), breakLoop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) stream.Iter[K, V, map[K]V] {
	h := m.Head()
	return stream.New(loop.Convert(h.Next, converter), loop.ToMap[K, V])
}

// Conv returns a breakable stream that applies the 'converter' function to the collection elements
func (m Map[K, V]) Conv(converter func(K, V) (K, V, error)) kvstream.Iter[K, V, map[K]V] {
	h := m.Head()
	return kvstream.New(breakLoop.Conv(breakLoop.From(h.Next), converter), breakLoop.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m Map[K, V]) Reduce(merge func(K, K, V, V) (K, V)) (K, V) {
	return map_.Reduce(m.elements, merge)
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	return map_.HasAny(m.elements, predicate)
}

func addToMap[K comparable, V any](key K, val V, order []K, uniques map[K]V) []K {
	if _, ok := uniques[key]; !ok {
		order = append(order, key)
		uniques[key] = val
	}
	return order
}
