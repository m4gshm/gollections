package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	oMapIter "github.com/m4gshm/gollections/immutable/ordered/map_/iter"
	"github.com/m4gshm/gollections/kv/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
	"github.com/m4gshm/gollections/slice/clone"
	sliceIter "github.com/m4gshm/gollections/slice/iter"
)

// NewMapKV converts a slice of key/value pairs to the Map.
func NewMapKV[K comparable, V any](elements []c.KV[K, V]) Map[K, V] {
	var (
		l       = len(elements)
		uniques = make(map[K]V, l)
		order   = make([]K, 0, l)
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		if _, ok := uniques[key]; !ok {
			order = append(order, key)
			uniques[key] = val
		}
	}
	return WrapMap(order, uniques)
}

// NewMap instantiates Map populated by the 'elements' map key/values
func NewMap[K comparable, V any](elements map[K]V, order []K) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for _, key := range order {
		uniques[key] = elements[key]
	}
	return WrapMap(clone.Of(order), uniques)
}

// WrapMap instantiates ordered Map using a map and an order slice as internal storage.
func WrapMap[K comparable, V any](order []K, elements map[K]V) Map[K, V] {
	return Map[K, V]{order: order, elements: elements, ksize: notsafe.GetTypeSize[K]()}
}

// Map is a collection that provides elements access by an unique key.
type Map[K comparable, V any] struct {
	order    []K
	elements map[K]V
	ksize    uintptr
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
func (m Map[K, V]) Head() oMapIter.OrderedMapIter[K, V] {
	return oMapIter.NewOrdered(m.elements, sliceIter.NewHeadS(m.order, m.ksize))
}

// First returns the first key/value pair of the map, an iterator to iterate over the remaining pair, and true\false marker of availability next pairs.
// If no more then ok==false.
func (m Map[K, V]) First() (oMapIter.OrderedMapIter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

// Tail creates an iterator pointing to the end of the map
func (m Map[K, V]) Tail() oMapIter.OrderedMapIter[K, V] {
	return oMapIter.NewOrdered(m.elements, sliceIter.NewTailS(m.order, m.ksize))
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
func (m Map[K, V]) Keys() c.Collection[K] {
	return m.K()
}

// K resutrns keys collection impl
func (m Map[K, V]) K() MapKeys[K] {
	return WrapKeys(m.order)
}

// Values resutrns values collection
func (m Map[K, V]) Values() c.Collection[V] {
	return m.V()
}

// V resutrns values collection impl
func (m Map[K, V]) V() MapValues[K, V] {
	return WrapVal(m.order, m.elements)
}

// Sort sorts the elements
func (m Map[K, V]) Sort(less slice.Less[K]) Map[K, V] {
	return m.sortBy(sort.Slice, less)
}

// StableSort sorts the elements
func (m Map[K, V]) StableSort(less slice.Less[K]) Map[K, V] {
	return m.sortBy(sort.SliceStable, less)
}

func (m Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) Map[K, V] {
	order := slice.Clone(m.order)
	slice.Sort(order, sorter, less)
	return WrapMap(order, m.elements)
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
func (m Map[K, V]) FilterKey(predicate func(K) bool) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Filter(h.Next, filter.Key[V](predicate)).Next, loop.ToMap[K, V])
}

// ConvertKey returns a stream that applies the 'converter' function to keys of the map
func (m Map[K, V]) ConvertKey(by func(K) K) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Convert(h.Next, convert.Key[V](by)).Next, loop.ToMap[K, V])
}

// FilterValue returns a stream consisting of key/value pairs where the value satisfies the condition of the 'predicate' function
func (m Map[K, V]) FilterValue(predicate func(V) bool) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Filter(h.Next, filter.Value[K](predicate)).Next, loop.ToMap[K, V])
}

// ConvertValue returns a stream that applies the 'converter' function to values of the map
func (m Map[K, V]) ConvertValue(by func(V) V) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Convert(h.Next, convert.Value[K](by)).Next, loop.ToMap[K, V])
}

// Filter returns a stream consisting of elements that satisfy the condition of the 'predicate' function
func (m Map[K, V]) Filter(predicate func(K, V) bool) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Filter(h.Next, predicate).Next, loop.ToMap[K, V])
}

// Convert returns a stream that applies the 'converter' function to the collection elements
func (m Map[K, V]) Convert(converter func(K, V) (K, V)) c.KVStream[K, V, map[K]V] {
	h := m.Head()
	return loop.Stream(loop.Convert(h.Next, converter).Next, loop.ToMap[K, V])
}

// Reduce reduces the key/value pairs of the map into an one pair using the 'merge' function
func (m Map[K, V]) Reduce(merge func(K, V, K, V) (K, V)) (K, V) {
	return map_.Reduce(m.elements, merge)
}

// HasAny finds the first key/value pair that satisfies the 'predicate' function condition and returns true if successful
func (m Map[K, V]) HasAny(predicate func(K, V) bool) bool {
	return map_.HasAny(m.elements, predicate)
}
