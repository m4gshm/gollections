package immutable

import (
	"fmt"

	breakLoop "github.com/m4gshm/gollections/break/kv/loop"
	breakMapFilter "github.com/m4gshm/gollections/break/kv/predicate"
	breakMapConvert "github.com/m4gshm/gollections/break/map_/convert"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/collection/immutable/ordered"
	"github.com/m4gshm/gollections/kv/convert"
	"github.com/m4gshm/gollections/kv/loop"
	filter "github.com/m4gshm/gollections/kv/predicate"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/slice"
)

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) Map[K, V] {
	return Map[K, V]{elements: elements}
}

// Map is a collection implementation that provides elements access by an unique key
type Map[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ collection.Map[int, any]                         = (*Map[int, any])(nil)
	_ collection.Map[int, any]                         = Map[int, any]{}
	_ c.KeyVal[MapKeys[int, any], MapValues[int, any]] = (*Map[int, any])(nil)
	_ c.KeyVal[MapKeys[int, any], MapValues[int, any]] = Map[int, any]{}
	_ fmt.Stringer                                     = (*Map[int, any])(nil)
	_ fmt.Stringer                                     = Map[int, any]{}
)

// All is used to iterate through the collection using `for ... range`. Supported since go 1.22 with GOEXPERIMENT=rangefunc enabled.
func (m Map[K, V]) All(consumer func(k K, v V) bool) {
	for k, v := range m.elements {
		if !consumer(k, v) {
			break
		}
	}
}

// Loop creates a loop to iterate through the collection.
func (m Map[K, V]) Loop() loop.Loop[K, V] {
	h := m.Head()
	return h.Next
}

// Deprecated: Head is deprecated. Will be replaced by rance-over function iterator.
// Head creates an iterator to iterate through the collection.
func (m Map[K, V]) Head() map_.Iter[K, V] {
	return map_.NewIter(m.elements)
}

// Deprecated: First is deprecated. Will be replaced by rance-over function iterator.
// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m Map[K, V]) First() (map_.Iter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Map collects the key/value pairs into a new map
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

// Sort returns sorted by keys map
func (m Map[K, V]) Sort(comparer slice.Comparer[K]) ordered.Map[K, V] {
	return m.sortBy(slice.Sort, comparer)
}

// StableSort returns sorted by keys map
func (m Map[K, V]) StableSort(comparer slice.Comparer[K]) ordered.Map[K, V] {
	return m.sortBy(slice.StableSort, comparer)
}

func (m Map[K, V]) sortBy(sorter func([]K, slice.Comparer[K]) []K, comparer slice.Comparer[K]) ordered.Map[K, V] {
	return ordered.WrapMap(sorter(map_.Keys(m.elements), comparer), m.elements)
}

func (m Map[K, V]) String() string {
	return map_.ToString(m.elements)
}

// Track applies the 'consumer' function for all key/value pairs until the consumer returns the c.Break to stop.
func (m Map[K, V]) Track(consumer func(K, V) error) error {
	return map_.Track(m.elements, consumer)
}

// TrackEach applies the 'consumer' function for every key/value pairs
func (m Map[K, V]) TrackEach(consumer func(K, V)) {
	map_.TrackEach(m.elements, consumer)
}

// FilterKey returns a loop consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterKey(predicate func(K) bool) loop.Loop[K, V] {
	return loop.Filter(m.Loop(), filter.Key[V](predicate))
}

// FiltKey returns a loop consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) breakLoop.Loop[K, V] {
	return loop.Filt(m.Loop(), breakMapFilter.Key[V](predicate))
}

// ConvertKey returns a loop that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(by func(K) K) loop.Loop[K, V] {
	return loop.Convert(m.Loop(), convert.Key[V](by))
}

// ConvKey returns a loop that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) breakLoop.Loop[K, V] {
	return loop.Conv(m.Loop(), breakMapConvert.Key[V](converter))
}

// FilterValue returns a loop consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) loop.Loop[K, V] {
	return loop.Filter(m.Loop(), filter.Value[K](predicate))
}

// FiltValue returns a loop consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltValue(predicate func(V) (bool, error)) breakLoop.Loop[K, V] {
	return loop.Filt(m.Loop(), breakMapFilter.Value[K](predicate))
}

// ConvertValue returns a loop that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(by func(V) V) loop.Loop[K, V] {
	return loop.Convert(m.Loop(), convert.Value[K](by))
}

// ConvValue returns a loop that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) breakLoop.Loop[K, V] {
	return loop.Conv(m.Loop(), breakMapConvert.Value[K](converter))
}

// Filter returns a loop consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) loop.Loop[K, V] {
	return loop.Filter(m.Loop(), predicate)
}

// Filt returns a breakable loop consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filt(predicate func(K, V) (bool, error)) breakLoop.Loop[K, V] {
	return loop.Filt(m.Loop(), predicate)
}

// Convert returns a loop that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) loop.Loop[K, V] {
	return loop.Convert(m.Loop(), converter)
}

// Conv returns a breakable loop that applies the 'converter' function to the collection elements
func (m Map[K, V]) Conv(converter func(K, V) (K, V, error)) breakLoop.Loop[K, V] {
	return loop.Conv(m.Loop(), converter)
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m Map[K, V]) Reduce(merge func(K, K, V, V) (K, V)) (K, V) {
	return map_.Reduce(m.elements, merge)
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	return map_.HasAny(m.elements, predicate)
}
