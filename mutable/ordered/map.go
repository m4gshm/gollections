package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/kviter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// AsMap converts a slice of key/value pairs to teh Map.
func AsMap[K comparable, V any](elements []c.KV[K, V]) *Map[K, V] {
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

// WrapMap instantiates ordered Map using a map and an order slice as internal storage.
func WrapMap[K comparable, V any](order []K, elements map[K]V) *Map[K, V] {
	return &Map[K, V]{order: order, elements: elements, ksize: notsafe.GetTypeSize[K]()}
}

// Map is the Collection implementation that provides element access by an unique key..
type Map[K comparable, V any] struct {
	order      []K
	elements   map[K]V
	changeMark int32
	ksize      uintptr
}

var (
	_ c.Settable[int, any]                                   = (*Map[int, any])(nil)
	_ c.SettableNew[int, any]                                = (*Map[int, any])(nil)
	_ c.SettableMap[int, any]                                = (*Map[int, any])(nil)
	_ c.ImmutableMapConvert[int, any, ordered.Map[int, any]] = (*Map[int, any])(nil)
	_ c.Map[int, any]                                        = (*Map[int, any])(nil)
	_ fmt.Stringer                                           = (*Map[int, any])(nil)
)

func (m *Map[K, V]) Begin() c.KVIterator[K, V] {
	h := m.Head()
	return &h
}

func (m *Map[K, V]) Head() iter.OrderedEmbedMapKVIter[K, V] {
	return iter.NewOrderedEmbedMapKV(m.elements, iter.NewHeadS(m.order, m.ksize))
}

func (m *Map[K, V]) Tail() iter.OrderedEmbedMapKVIter[K, V] {
	return iter.NewOrderedEmbedMapKV(m.elements, iter.NewTailS(m.order, m.ksize))
}

func (m *Map[K, V]) First() (iter.OrderedEmbedMapKVIter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

func (m *Map[K, V]) Collect() map[K]V {
	return map_.Clone(m.elements)
}

func (m *Map[K, V]) Map() map[K]V {
	return m.Collect()
}

// Sort transforms to the ordered Map contains sorted elements.
func (m *Map[K, V]) Sort(less slice.Less[K]) *Map[K, V] {
	return m.sortBy(sort.Slice, less)
}

func (m *Map[K, V]) StableSort(less slice.Less[K]) *Map[K, V] {
	return m.sortBy(sort.SliceStable, less)
}

func (m *Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) *Map[K, V] {
	slice.Sort(m.order, sorter, less)
	return m
}

func (m *Map[K, V]) Len() int {
	return len(m.order)
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return map_.ForOrdered(m.order, m.elements, walker)
}

func (m *Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	map_.ForEachOrdered(m.order, m.elements, walker)
}

func (m *Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.TrackOrdered(m.order, m.elements, tracker)
}

func (m *Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEachOrdered(m.order, m.elements, tracker)
}

func (m *Map[K, V]) Contains(key K) bool {
	_, ok := m.elements[key]
	return ok
}

func (m *Map[K, V]) Get(key K) (V, bool) {
	val, ok := m.elements[key]
	return val, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	m.SetNew(key, value)
}

func (m *Map[K, V]) SetNew(key K, value V) bool {
	var (
		u     = m.elements
		_, ok = u[key]
	)
	if !ok {
		u[key] = value
		m.order = append(m.order, key)
	}
	return !ok
}

func (m *Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return m.K()
}

func (m *Map[K, V]) K() ordered.MapKeys[K] {
	return ordered.WrapKeys(m.order)
}

func (m *Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return m.V()
}

func (m *Map[K, V]) V() ordered.MapValues[K, V] {
	return ordered.WrapVal(m.order, m.elements)
}

func (m *Map[K, V]) String() string {
	return map_.ToStringOrdered(m.order, m.elements)
}

func (m *Map[K, V]) FilterKey(filter func(K) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, func(key K, val V) bool { return filter(key) }), kviter.ToMap[K, V])
}

func (m *Map[K, V]) ConvertKey(by func(K) K) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, func(key K, val V) (K, V) { return by(key), val }), kviter.ToMap[K, V])
}

func (m *Map[K, V]) FilterValue(filter func(V) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, func(key K, val V) bool { return filter(val) }), kviter.ToMap[K, V])
}

func (m *Map[K, V]) ConvertValue(by func(V) V) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, func(key K, val V) (K, V) { return key, by(val) }), kviter.ToMap[K, V])
}

func (m *Map[K, V]) Filter(filter func(K, V) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, filter), kviter.ToMap[K, V])
}

func (m *Map[K, V]) Convert(by func(K, V) (K, V)) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, by), kviter.ToMap[K, V])
}

func (m *Map[K, V]) Reduce(by c.Quaternary[K, V]) (K, V) {
	h := m.Head()
	return loop.ReduceKV(h.Next, by)
}

func (m *Map[K, V]) Immutable() ordered.Map[K, V] {
	return ordered.WrapMap(slice.Clone(m.order), map_.Clone(m.elements))
}

func (m *Map[K, V]) SetMap(kvs c.Map[K, V]) {
	kvs.TrackEach(func(key K, value V) { m.Set(key, value) })
}
