package ordered

import (
	"fmt"
	"sort"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/kvit"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/notsafe"
	"github.com/m4gshm/gollections/slice"
)

// ConvertKVsToMap converts a slice of key/value pairs to the Map.
func ConvertKVsToMap[K comparable, V any](elements []c.KV[K, V]) Map[K, V] {
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
func NewMap[K comparable, V any](elements map[K]V) Map[K, V] {
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
func WrapMap[K comparable, V any](order []K, elements map[K]V) Map[K, V] {
	return Map[K, V]{order: order, elements: elements, ksize: notsafe.GetTypeSize[K]()}
}

// Map is the Collection implementation that provides element access by an unique key.
type Map[K comparable, V any] struct {
	order    []K
	elements map[K]V
	ksize    uintptr
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ c.Map[int, any] = Map[int, any]{}
	_ fmt.Stringer    = (*Map[int, any])(nil)
	_ fmt.Stringer    = Map[int, any]{}
)

func (s Map[K, V]) Begin() c.KVIterator[K, V] {
	h := s.Head()
	return &h
}

func (s Map[K, V]) Head() it.OrderedEmbedMapKVIter[K, V] {
	return it.NewOrderedEmbedMapKV(s.elements, it.NewHeadS(s.order, s.ksize))
}

func (s Map[K, V]) Tail() it.OrderedEmbedMapKVIter[K, V] {
	return it.NewOrderedEmbedMapKV(s.elements, it.NewTailS(s.order, s.ksize))
}

func (s Map[K, V]) First() (it.OrderedEmbedMapKVIter[K, V], K, V, bool) {
	var (
		iter               = s.Head()
		firstK, firstV, ok = iter.Next()
	)
	return iter, firstK, firstV, ok
}

func (s Map[K, V]) Collect() map[K]V {
	return map_.Clone(s.elements)
}

func (s Map[K, V]) Map() map[K]V {
	return s.Collect()
}

func (s Map[K, V]) Len() int {
	return len(s.order)
}

func (s Map[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Map[K, V]) Contains(key K) bool {
	_, ok := s.elements[key]
	return ok
}

func (s Map[K, V]) Get(key K) (V, bool) {
	val, ok := s.elements[key]
	return val, ok
}

func (s Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return s.K()
}

func (s Map[K, V]) K() MapKeys[K] {
	return WrapKeys(s.order)
}

func (s Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return s.V()
}

func (s Map[K, V]) V() MapValues[K, V] {
	return WrapVal(s.order, s.elements)
}

func (s Map[K, V]) Sort(less slice.Less[K]) Map[K, V] {
	return s.sortBy(sort.Slice, less)
}

func (s Map[K, V]) StableSort(less slice.Less[K]) Map[K, V] {
	return s.sortBy(sort.SliceStable, less)
}

func (s Map[K, V]) sortBy(sorter slice.Sorter, less slice.Less[K]) Map[K, V] {
	c := slice.Clone(s.order)
	slice.Sort(c, sorter, less)
	return WrapMap(c, s.elements)
}

func (s Map[K, V]) String() string {
	return map_.ToStringOrdered(s.order, s.elements)
}

func (s Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.TrackOrdered(s.order, s.elements, tracker)
}

func (s Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEachOrdered(s.order, s.elements, tracker)
}

func (s Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return map_.ForOrdered(s.order, s.elements, walker)
}

func (s Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	map_.ForEachOrdered(s.order, s.elements, walker)
}

func (s Map[K, V]) FilterKey(filter func(K) bool) c.MapPipe[K, V, map[K]V] {
	h := s.Head()
	return it.NewKVPipe(it.FilterKV(&h, func(key K, val V) bool { return filter(key) }), kvit.ToMap[K, V])
}

func (s Map[K, V]) ConvertKey(by func(K) K) c.MapPipe[K, V, map[K]V] {
	h := s.Head()
	return it.NewKVPipe(it.ConvertKV(&h, func(key K, val V) (K, V) { return by(key), val }), kvit.ToMap[K, V])
}

func (s Map[K, V]) FilterValue(filter func(V) bool) c.MapPipe[K, V, map[K]V] {
	h := s.Head()
	return it.NewKVPipe(it.FilterKV(&h, func(key K, val V) bool { return filter(val) }), kvit.ToMap[K, V])
}

func (s Map[K, V]) ConvertValue(by func(V) V) c.MapPipe[K, V, map[K]V] {
	h := s.Head()
	return it.NewKVPipe(it.ConvertKV(&h, func(key K, val V) (K, V) { return key, by(val) }), kvit.ToMap[K, V])
}

func (s Map[K, V]) Filter(filter func(K, V) bool) c.MapPipe[K, V, map[K]V] {
	h := s.Head()
	return it.NewKVPipe(it.FilterKV(&h, filter), kvit.ToMap[K, V])
}

func (s Map[K, V]) Convert(by func(K, V) (K, V)) c.MapPipe[K, V, map[K]V] {
	h := s.Head()
	return it.NewKVPipe(it.ConvertKV(&h, by), kvit.ToMap[K, V])
}

func (s Map[K, V]) Reduce(by c.Quaternary[K, V]) (K, V) {
	h := s.Head()
	return loop.ReduceKV(h.Next, by)
}
