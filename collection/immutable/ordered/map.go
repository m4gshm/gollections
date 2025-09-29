package ordered

import (
	"fmt"

	converte "github.com/m4gshm/gollections/break/kv/convert"
	filtere "github.com/m4gshm/gollections/break/kv/predicate"
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collection"
	"github.com/m4gshm/gollections/kv/convert"
	filter "github.com/m4gshm/gollections/kv/predicate"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/seq2"
	"github.com/m4gshm/gollections/slice"
)

// WrapMap instantiates an ordered Map using a map and an order slice as internal storage.
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

func (m Map[K, V]) Head() (K, V, bool) {
	return seq2.Head(m.All)
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

// FilterKey returns a seq consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterKey(predicate func(K) bool) collection.Seq2[K, V] {
	return seq2.Filter(m.All, filter.Key[V](predicate))
}

// FiltKey returns an errorable seq consisting of key/value pairs where the key satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltKey(predicate func(K) (bool, error)) collection.SeqE[c.KV[K, V]] {
	return seq2.Filt(m.All, filtere.Key[V](predicate))
}

// ConvertKey returns a seq that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(converter func(K) K) collection.Seq2[K, V] {
	return seq2.Convert(m.All, convert.Key[V](converter))
}

// ConvKey returns an errorable seq that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvKey(converter func(K) (K, error)) collection.SeqE[c.KV[K, V]] {
	return seq2.Conv(m.All, converte.Key[V](converter))
}

// FilterValue returns a seq consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) collection.Seq2[K, V] {
	return seq2.Filter(m.All, filter.Value[K](predicate))
}

// FiltValue returns an errorable seq consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FiltValue(predicate func(V) (bool, error)) collection.SeqE[c.KV[K, V]] {
	return seq2.Filt(m.All, filtere.Value[K](predicate))
}

// ConvertValue returns a seq that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(converter func(V) V) collection.Seq2[K, V] {
	return seq2.Convert(m.All, convert.Value[K](converter))
}

// ConvValue returns a errorable seq that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvValue(converter func(V) (V, error)) collection.SeqE[c.KV[K, V]] {
	return seq2.Conv(m.All, converte.Value[K](converter))
}

// Filter returns a seq consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) collection.Seq2[K, V] {
	return seq2.Filter(m.All, predicate)
}

// Filt returns a errorable seq consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filt(predicate func(K, V) (bool, error)) collection.SeqE[c.KV[K, V]] {
	return seq2.Filt(m.All, predicate)
}

// Convert returns a seq that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) collection.Seq2[K, V] {
	return seq2.Convert(m.All, converter)
}

// Conv returns a errorable seq that applies the 'converter' function to the collection elements
func (m Map[K, V]) Conv(converter func(K, V) (K, V, error)) collection.SeqE[c.KV[K, V]] {
	return seq2.Conv(m.All, converter)
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
