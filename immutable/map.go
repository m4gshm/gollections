package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	kviter "github.com/m4gshm/gollections/loop/kv/iter"
	"github.com/m4gshm/gollections/loop/kv/stream"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/map_/iter"
	"github.com/m4gshm/gollections/slice"
)

// ConvertKVsToMap converts a slice of key/value pairs to the Map.
func ConvertKVsToMap[K comparable, V any](elements []c.KV[K, V]) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, kv := range elements {
		uniques[kv.Key()] = kv.Value()
	}
	return WrapMap(uniques)
}

// NewMap instantiates Map populated by the 'elements' map key/values
func NewMap[K comparable, V any](elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) Map[K, V] {
	return Map[K, V]{elements: elements}
}

// Map is a collection that provides elements access by an unique key.
type Map[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ fmt.Stringer    = (*Map[int, any])(nil)
	_ c.Map[int, any] = Map[int, any]{}
	_ fmt.Stringer    = Map[int, any]{}
)

// Begin creates iterator
func (m Map[K, V]) Begin() c.KVIterator[K, V] {
	h := m.Head()
	return &h
}

// Head creates iterator
func (m Map[K, V]) Head() iter.EmbedMapKVIter[K, V] {
	return iter.New(m.elements)
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m Map[K, V]) First() (iter.EmbedMapKVIter[K, V], K, V, bool) {
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
func (m Map[K, V]) Keys() c.Collection[K] {
	return m.K()
}

// K resutrns keys collection impl
func (m Map[K, V]) K() MapKeys[K, V] {
	return WrapKeys(m.elements)
}

// Values resutrns values collection
func (m Map[K, V]) Values() c.Collection[V] {
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

// FilterKey returns a pipe consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterKey(predicate func(K) bool) c.MapStream[K, V, map[K]V] {
	h := m.Head()
	return stream.New(kviter.Filter(h.Next, filter.Key[V](predicate)).Next, stream.ToMap[K, V])
}

// ConvertKey returns a pipe that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(by func(K) K) c.MapStream[K, V, map[K]V] {
	h := m.Head()
	return stream.New(kviter.Convert(h.Next, convert.Key[V](by)).Next, stream.ToMap[K, V])
}

// FilterValue returns a pipe consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) c.MapStream[K, V, map[K]V] {
	h := m.Head()
	return stream.New(kviter.Filter(h.Next, filter.Value[K](predicate)).Next, stream.ToMap[K, V])
}

// ConvertValue returns a pipe that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(by func(V) V) c.MapStream[K, V, map[K]V] {
	h := m.Head()
	return stream.New(kviter.Convert(h.Next, convert.Value[K](by)).Next, stream.ToMap[K, V])
}

// Filter returns a pipe consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) c.MapStream[K, V, map[K]V] {
	h := m.Head()
	return stream.New(kviter.Filter(h.Next, predicate).Next, stream.ToMap[K, V])
}

// Convert returns a pipe that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) c.MapStream[K, V, map[K]V] {
	h := m.Head()
	return stream.New(kviter.Convert(h.Next, converter).Next, stream.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m Map[K, V]) Reduce(merge func(K, V, K, V) (K, V)) (K, V) {
	return map_.Reduce(m.elements, merge)
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	return map_.HasAny(m.elements, predicate)
}
