package ordered

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/kv/loop"
	breakMapFilter "github.com/m4gshm/gollections/break/kv/predicate"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/kv/convert"
	kvloop "github.com/m4gshm/gollections/kv/loop"
	filter "github.com/m4gshm/gollections/kv/predicate"
	"github.com/m4gshm/gollections/map_"
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

// All is used to iterate through the collection using `for key, val := range`.
func (m Map[K, V]) All(consumer func(K, V) bool) {
	map_.TrackOrderedWhile(m.order, m.elements, consumer)
}

// Loop creates a loop to iterate through the collection.
// Deprecated: replaced by the [All].
func (m Map[K, V]) Loop() kvloop.Loop[K, V] {
	h := m.Head()
	return h.Next
}

// Head creates an iterator to iterate through the collection.
// Deprecated: replaced by the [All].
func (m Map[K, V]) Head() MapIter[K, V] {
	return NewMapIter(m.elements, slice.NewHead(m.order))
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
// Deprecated: replaced by the [All].
func (m Map[K, V]) First() (MapIter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Tail creates an iterator pointing to the end of the map
// Deprecated: Tail is deprecated. Will be replaced by a rance-over function iterator.
func (m Map[K, V]) Tail() MapIter[K, V] {
	return NewMapIter(m.elements, slice.NewTail(m.order))
}

// Map collects the key/value pairs into a new map
func (m Map[K, V]) Map() map[K]V {
	return map_.Clone(m.elements)
}

// Len returns amount of elements
func (m Map[K, V]) Len() int {
	return len(m.order)
}

// IsEmpty returns true if the map is empty
func (m Map[K, V]) IsEmpty() bool {
	return collection.IsEmpty(m)
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

// Track applies the 'consumer' function for all key/value pairs until the consumer returns the c.Break to stop.
func (m Map[K, V]) Track(consumer func(K, V) error) error {
	return map_.TrackOrdered(m.order, m.elements, consumer)
}

// TrackEach applies the 'consumer' function for every key/value pairs
func (m Map[K, V]) TrackEach(consumer func(K, V)) {
	map_.TrackEachOrdered(m.order, m.elements, consumer)
}

// FilterKey returns a loop consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterKey(predicate func(K) bool) kvloop.Loop[K, V] {
	h := m.Head()
	return kvloop.Filter(h.Next, filter.Key[V](predicate))
}

// FiltKey returns a loop consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) breakLoop.Loop[K, V] {
	h := m.Head()
	return kvloop.Filt(h.Next, breakMapFilter.Key[V](predicate))
}

// ConvertKey returns a loop that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(by func(K) K) kvloop.Loop[K, V] {
	h := m.Head()
	return kvloop.Convert(h.Next, convert.Key[V](by))
}

// ConvKey returns a loop that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) breakLoop.Loop[K, V] {
	h := m.Head()
	return kvloop.Conv(h.Next, breakMapConvert.Key[V](converter))
}

// FilterValue returns a loop consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) kvloop.Loop[K, V] {
	h := m.Head()
	return kvloop.Filter(h.Next, filter.Value[K](predicate))
}

// FiltValue returns a breakable loop consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltValue(predicate func(V) (bool, error)) breakLoop.Loop[K, V] {
	h := m.Head()
	return kvloop.Filt(h.Next, breakMapFilter.Value[K](predicate))
}

// ConvertValue returns a loop that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(converter func(V) V) kvloop.Loop[K, V] {
	h := m.Head()
	return kvloop.Convert(h.Next, convert.Value[K](converter))
}

// ConvValue returns a breakable loop that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) breakLoop.Loop[K, V] {
	h := m.Head()
	return kvloop.Conv(h.Next, breakMapConvert.Value[K](converter))
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) kvloop.Loop[K, V] {
	h := m.Head()
	return kvloop.Filter(h.Next, predicate)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filt(predicate func(K, V) (bool, error)) breakLoop.Loop[K, V] {
	h := m.Head()
	return kvloop.Filt(h.Next, predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) kvloop.Loop[K, V] {
	h := m.Head()
	return kvloop.Convert(h.Next, converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (m Map[K, V]) Conv(converter func(K, V) (K, V, error)) breakLoop.Loop[K, V] {
	h := m.Head()
	return kvloop.Conv(h.Next, converter)
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
