package ordered

import (
	"fmt"

	converte "github.com/m4gshm/gollections/break/kv/convert"
	filtere "github.com/m4gshm/gollections/break/kv/predicate"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/kv/convert"
	filter "github.com/m4gshm/gollections/kv/predicate"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq"
	"github.com/m4gshm/gollections/seq2"
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
	_ c.SettableMap[c.TrackEach[int, any]]                        = (*Map[int, any])(nil)
	_ c.ImmutableMapConvert[ordered.Map[int, any]]                = (*Map[int, any])(nil)
	_ collection.Map[int, any]                                    = (*Map[int, any])(nil)
	_ c.KeyVal[ordered.MapKeys[int], ordered.MapValues[int, any]] = (*Map[int, any])(nil)
	_ fmt.Stringer                                                = (*Map[int, any])(nil)
)

// All is used to iterate through the collection using `for key, val := range`.
func (m *Map[K, V]) All(consumer func(K, V) bool) {
	if m != nil {
		map_.TrackOrderedWhile(m.order, m.elements, consumer)
	}
}

func (m *Map[K, V]) Head() (K, V, bool) {
	return seq2.Head(m.All)
}

// Map collects the key/value pairs into a new map
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
	return collection.IsEmpty(m)
}

// Track applies the 'consumer' function for all key/value pairs until the consumer returns the c.Break to stop.
func (m *Map[K, V]) Track(consumer func(K, V) error) error {
	if m == nil {
		return nil
	}
	return map_.TrackOrdered(m.order, m.elements, consumer)
}

// TrackEach applies the 'consumer' function for every key/value pairs
func (m *Map[K, V]) TrackEach(consumer func(K, V)) {
	if m == nil {
		return
	}
	map_.TrackEachOrdered(m.order, m.elements, consumer)
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

// FilterKey returns a seq consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterKey(predicate func(K) bool) seq.Seq2[K, V] {
	return seq2.Filter(m.All, filter.Key[V](predicate))
}

// FiltKey returns a seq consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FiltKey(predicate func(K) (bool, error)) seq.SeqE[c.KV[K, V]] {
	return seq2.Filt(m.All, filtere.Key[V](predicate))
}

// ConvertKey returns a seq that applies the 'converter' function to keys of the map
func (m *Map[K, V]) ConvertKey(converter func(K) K) seq.Seq2[K, V] {
	return seq2.Convert(m.All, convert.Key[V](converter))
}

// ConvKey returns a seq that applies the 'converter' function to keys of the map
func (m *Map[K, V]) ConvKey(converter func(K) (K, error)) seq.SeqE[c.KV[K, V]] {
	return seq2.Conv(m.All, converte.Key[V](converter))
}

// FilterValue returns a seq consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FilterValue(predicate func(V) bool) seq.Seq2[K, V] {
	return seq2.Filter(m.All, filter.Value[K](predicate))
}

// FiltValue returns a errorable seq consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m *Map[K, V]) FiltValue(predicate func(V) (bool, error)) seq.SeqE[c.KV[K, V]] {
	return seq2.Filt(m.All, filtere.Value[K](predicate))
}

// ConvertValue returns a seq that applies the 'converter' function to values of the map
func (m *Map[K, V]) ConvertValue(converter func(V) V) seq.Seq2[K, V] {
	return seq2.Convert(m.All, convert.Value[K](converter))
}

// ConvValue returns a errorable seq that applies the 'converter' function to values of the map
func (m *Map[K, V]) ConvValue(converter func(V) (V, error)) seq.SeqE[c.KV[K, V]] {
	return seq2.Conv(m.All, converte.Value[K](converter))
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filter(predicate func(K, V) bool) seq.Seq2[K, V] {
	return seq2.Filter(m.All, predicate)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'predicate' function
func (m *Map[K, V]) Filt(predicate func(K, V) (bool, error)) seq.SeqE[c.KV[K, V]] {
	return seq2.Filt(m.All, predicate)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Convert(converter func(K, V) (K, V)) seq.Seq2[K, V] {
	return seq2.Convert(m.All, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (m *Map[K, V]) Conv(converter func(K, V) (K, V, error)) seq.SeqE[c.KV[K, V]] {
	return seq2.Conv(m.All, converter)
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
func (m *Map[K, V]) SetMap(kvs c.TrackEach[K, V]) {
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
