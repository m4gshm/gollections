package immutable

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/kviter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/filter"
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

// NewMap instantiates Map with values copied from an map.
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

// Map is the Collection implementation that provides element access by an unique key.
type Map[K comparable, V any] struct {
	elements map[K]V
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ c.Map[int, any] = Map[int, any]{}
	_ fmt.Stringer    = (*Map[int, any])(nil)
	_ fmt.Stringer    = Map[int, any]{}
)

// Begin creates a key/value iterator interface.
func (m Map[K, V]) Begin() c.KVIterator[K, V] {
	h := m.Head()
	return &h
}

// Head creates a key/value iterator implementation started from the head.
func (m Map[K, V]) Head() iter.EmbedMapKVIter[K, V] {
	return iter.NewEmbedMapKV(m.elements)
}

// Map exports the content as a map.
func (m Map[K, V]) Map() map[K]V {
	return map_.Clone(m.elements)
}

// Sort transforms to the ordered Map contains sorted elements.
func (m Map[K, V]) Sort(less slice.Less[K]) ordered.Map[K, V] {
	return m.sortBy(sort.Slice, less)
}

func (m Map[K, V]) StableSort(less slice.Less[K]) ordered.Map[K, V] {
	return m.sortBy(sort.SliceStable, less)
}

func (m Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) ordered.Map[K, V] {
	c := map_.Keys(m.elements)
	slice.Sort(c, sorter, less)
	return ordered.WrapMap(c, m.elements)
}

// String is part of the Stringer interface for printing the string representation of this Map.
func (m Map[K, V]) String() string {
	return map_.ToString(m.elements)
}

// Track apply a tracker to touch key, value from the inside. To stop traking just return the map_.Break.
func (m Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.Track(m.elements, tracker)
}

// TrackEach apply a tracker to touch each key, value from the inside.
func (m Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEach(m.elements, tracker)
}

func (m Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return map_.For(m.elements, walker)
}

func (m Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	map_.ForEach(m.elements, walker)
}

func (m Map[K, V]) Len() int {
	return len(m.elements)
}

func (m Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m Map[K, V]) Contains(key K) bool {
	_, ok := m.elements[key]
	return ok
}

func (m Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.elements[key]
	return val, ok
}

func (m Map[K, V]) Keys() c.Collection[K] {
	return m.K()
}

func (m Map[K, V]) K() MapKeys[K, V] {
	return WrapKeys(m.elements)
}

func (m Map[K, V]) Values() c.Collection[V] {
	return m.V()
}

func (m Map[K, V]) V() MapValues[K, V] {
	return WrapVal(m.elements)
}

func (m Map[K, V]) FilterKey(predicate func(K) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, filter.Key[K, V](predicate)), kviter.ToMap[K, V])
}

func (m Map[K, V]) ConvertKey(by func(K) K) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, func(key K, val V) (K, V) { return by(key), val }), kviter.ToMap[K, V])
}

func (m Map[K, V]) FilterValue(predicate func(V) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, filter.Value[K](predicate)), kviter.ToMap[K, V])
}

func (m Map[K, V]) ConvertValue(by func(V) V) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, func(key K, val V) (K, V) { return key, by(val) }), kviter.ToMap[K, V])
}

func (m Map[K, V]) Filter(filter func(K, V) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, filter), kviter.ToMap[K, V])
}

// Map creates an Iterator that applies a converter to iterable elements.
// The 'by' converter is applied only when the 'Next' method of returned Iterator is called and does not change the elements of the map.
func (m Map[K, V]) Convert(by func(K, V) (K, V)) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, by), kviter.ToMap[K, V])
}

// Reduce reduces key\value pairs to an one.
func (m Map[K, V]) Reduce(by c.Quaternary[K, V]) (K, V) {
	h := m.Head()
	return loop.ReduceKV(h.Next, by)
}
