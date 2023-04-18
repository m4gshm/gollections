package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/kviter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// ConvertKVsToMap converts a slice of key/value pairs to the Map.
func ConvertKVsToMap[K comparable, V any](elements []c.KV[K, V]) *Map[K, V] {
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

// NewMap instantiates Map and copies elements to it.
func NewMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	var (
		uniques = make(map[K]V, len(elements))
		order   = make([]K, len(elements))
	)
	for key, val := range elements {
		order = append(order, key)
		uniques[key] = val
	}
	return WrapMap(order, uniques)
}

// WrapMap instantiates ordered Map using a map and an order slice as internal storage.
func WrapMap[K comparable, V any](order []K, elements map[K]V) *Map[K, V] {
	return &Map[K, V]{order: order, elements: elements, ksize: notsafe.GetTypeSize[K]()}
}

// Map is the Collection implementation that provides element access by an unique key.
type Map[K comparable, V any] struct {
	order    []K
	elements map[K]V
	ksize    uintptr
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ fmt.Stringer    = (*Map[int, any])(nil)
)

func (m *Map[K, V]) Begin() c.KVIterator[K, V] {
	h := m.Head()
	return &h
}

func (m *Map[K, V]) Head() iter.OrderedEmbedMapKVIter[K, V] {
	var (
		order    []K
		elements map[K]V
		ksize    uintptr
	)
	if m != nil {
		order = m.order
		elements = m.elements
		ksize = m.ksize
	}

	return *iter.NewOrderedEmbedMapKV(elements, iter.NewHeadS(order, ksize))
}

func (m *Map[K, V]) Tail() iter.OrderedEmbedMapKVIter[K, V] {
	var (
		order    []K
		elements map[K]V
		ksize    uintptr
	)
	if m != nil {
		order = m.order
		elements = m.elements
		ksize = m.ksize
	}
	return *iter.NewOrderedEmbedMapKV(elements, iter.NewTailS(order, ksize))
}

func (m *Map[K, V]) First() (iter.OrderedEmbedMapKVIter[K, V], K, V, bool) {
	var (
		iterator           = m.Head()
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

func (m *Map[K, V]) Map() map[K]V {
	var elements map[K]V
	if m != nil {
		elements = m.elements
	}
	return map_.Clone(elements)
}

func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.order)
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Map[K, V]) Contains(key K) bool {
	if m == nil || m.elements == nil {
		return false
	}
	_, ok := m.elements[key]
	return ok
}

func (m *Map[K, V]) Get(key K) (val V, ok bool) {
	if m == nil || m.elements == nil {
		return
	}
	val, ok = m.elements[key]
	return val, ok
}

func (m *Map[K, V]) Keys() c.Collection[K] {
	return m.K()
}

func (m *Map[K, V]) K() MapKeys[K] {
	var order []K
	if m != nil {
		order = m.order
	}
	return WrapKeys(order)
}

func (m *Map[K, V]) Values() c.Collection[V] {
	return m.V()
}

func (m *Map[K, V]) V() MapValues[K, V] {
	var order []K
	var elements map[K]V
	if m != nil {
		order = m.order
		elements = m.elements
	}
	return WrapVal(order, elements)
}

func (m *Map[K, V]) Sort(less slice.Less[K]) *Map[K, V] {
	return m.sortBy(sort.Slice, less)
}

func (m *Map[K, V]) StableSort(less slice.Less[K]) *Map[K, V] {
	return m.sortBy(sort.SliceStable, less)
}

func (m *Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) *Map[K, V] {
	var order []K
	var elements map[K]V
	if m != nil {
		order = slice.Clone(m.order)
		elements = m.elements
	}
	slice.Sort(order, sorter, less)
	return WrapMap(order, elements)
}

func (m *Map[K, V]) String() string {
	if m == nil {
		return ""
	}
	return map_.ToStringOrdered(m.order, m.elements)
}

func (m *Map[K, V]) Track(tracker func(K, V) error) error {
	if m == nil {
		return nil
	}
	return map_.TrackOrdered(m.order, m.elements, tracker)
}

func (m *Map[K, V]) TrackEach(tracker func(K, V)) {
	if m != nil {
		map_.TrackEachOrdered(m.order, m.elements, tracker)
	}
}

func (m *Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	if m == nil {
		return nil
	}
	return map_.ForOrdered(m.order, m.elements, walker)
}

func (m *Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	if m != nil {
		map_.ForEachOrdered(m.order, m.elements, walker)
	}
}

func (m *Map[K, V]) FilterKey(predicate func(K) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, filter.Key[V](predicate)), kviter.ToMap[K, V])
}

func (m *Map[K, V]) ConvertKey(by func(K) K) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, convert.Key[V](by)), kviter.ToMap[K, V])
}

func (m *Map[K, V]) FilterValue(predicate func(V) bool) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.FilterKV(&h, filter.Value[K](predicate)), kviter.ToMap[K, V])
}

func (m *Map[K, V]) ConvertValue(by func(V) V) c.MapPipe[K, V, map[K]V] {
	h := m.Head()
	return iter.NewKVPipe(iter.ConvertKV(&h, convert.Value[K](by)), kviter.ToMap[K, V])
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
