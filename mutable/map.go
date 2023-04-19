package mutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/iter/impl/iter"
	"github.com/m4gshm/gollections/kviter"
	"github.com/m4gshm/gollections/loop"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/map_/convert"
	"github.com/m4gshm/gollections/map_/filter"
)

// NewMap instantiates Map with a predefined capacity.
func NewMap[K comparable, V any](capacity int) *Map[K, V] {
	return WrapMap(make(map[K]V, capacity))
}

// AsMap converts a slice of key/value pairs into a Map instance.
func AsMap[K comparable, V any](elements []c.KV[K, V]) *Map[K, V] {
	return WrapMap(map_.Of(elements...))
}

// ToMap instantiates Map and copies elements to it.
func ToMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

// WrapMap instantiates Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	m := Map[K, V](elements)
	return &m
}

// Map is the Collection implementation based on the embedded map.
type Map[K comparable, V any] map[K]V

var (
	_ c.Deleteable[int] = (*Map[int, any])(nil)
	// _ c.Deleteable[int]                                        = (Map[int, any])(nil)
	_ c.Removable[int, any] = (*Map[int, any])(nil)
	// _ c.Removable[int, any]                                    = (Map[int, any])(nil)
	_ c.Settable[int, any] = (*Map[int, any])(nil)
	// _ c.Settable[int, any]                                     = (Map[int, any])(nil)
	_ c.SettableNew[int, any] = (*Map[int, any])(nil)
	// _ c.SettableNew[int, any]                                  = (Map[int, any])(nil)
	_ c.SettableMap[int, any] = (*Map[int, any])(nil)
	// _ c.SettableMap[int, any]                                  = (Map[int, any])(nil)
	_ c.ImmutableMapConvert[int, any, *immutable.Map[int, any]] = (*Map[int, any])(nil)
	// _ c.ImmutableMapConvert[int, any, immutable.Map[int, any]] = (Map[int, any])(nil)
	_ c.Map[int, any] = (*Map[int, any])(nil)
	// _ c.Map[int, any]                                          = (Map[int, any])(nil)
	_ fmt.Stringer = (*Map[int, any])(nil)
	// _ fmt.Stringer                                             = (Map[int, any])(nil)
)

func (m *Map[K, V]) Begin() c.KVIterator[K, V] {
	h := m.Head()
	return &h
}

func (m *Map[K, V]) Head() iter.EmbedMapKVIter[K, V] {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return *iter.NewEmbedMapKV(out)
}

func (m *Map[K, V]) First() (iter.EmbedMapKVIter[K, V], K, V, bool) {
	var out map[K]V
	if m != nil {
		out = *m
	}
	var (
		iterator           = *iter.NewEmbedMapKV(out)
		firstK, firstV, ok = iterator.Next()
	)
	return iterator, firstK, firstV, ok
}

func (m *Map[K, V]) Map() (out map[K]V) {
	if m == nil {
		return
	}
	return map_.Clone(*m)
}

func (m *Map[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(*m)
}

func (m *Map[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

func (m *Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	if m == nil {
		return nil
	}
	return map_.For(*m, walker)
}

func (m *Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	if m != nil {
		map_.ForEach(*m, walker)
	}
}

func (m *Map[K, V]) Track(tracker func(K, V) error) error {
	if m == nil {
		return nil
	}
	return map_.Track(*m, tracker)
}

func (m *Map[K, V]) TrackEach(tracker func(K, V)) {
	if m != nil {
		map_.TrackEach(*m, tracker)
	}
}

func (m *Map[K, V]) Contains(key K) (ok bool) {
	if m != nil {
		_, ok = (*m)[key]
	}
	return ok
}

func (m *Map[K, V]) Get(key K) (val V, ok bool) {
	if m != nil {
		val, ok = (*m)[key]
	}
	return val, ok
}

func (m *Map[K, V]) Set(key K, value V) {
	if m == nil {
		return
	} else if (*m) == nil {
		*m = Map[K, V]{}
	}
	(*m)[key] = value
}

func (m *Map[K, V]) SetNew(key K, value V) bool {
	if m == nil {
		return false
	} else if (*m) == nil {
		*m = Map[K, V]{}
	}

	if _, ok := (*m)[key]; !ok {
		(*m)[key] = value
		return true
	}
	return false
}

func (m *Map[K, V]) Delete(keys ...K) {
	for _, key := range keys {
		m.DeleteOne(key)
	}
}

func (m *Map[K, V]) DeleteOne(key K) {
	if m != nil {
		delete(*m, key)
	}
}

func (m *Map[K, V]) Remove(key K) (v V, ok bool) {
	if m == nil {
		return v, ok
	}
	v, ok = m.Get(key)
	m.Delete(key)
	return v, ok
}

func (m *Map[K, V]) Keys() c.Collection[K] {
	return m.K()
}

func (m *Map[K, V]) K() *immutable.MapKeys[K, V] {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return immutable.WrapKeys(out)
}

func (m *Map[K, V]) Values() c.Collection[V] {
	return m.V()
}

func (m *Map[K, V]) V() *immutable.MapValues[K, V] {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return immutable.WrapVal(out)
}

func (m *Map[K, V]) String() string {
	var out map[K]V
	if m != nil {
		out = *m
	}
	return map_.ToString(out)
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

func (m *Map[K, V]) Immutable() *immutable.Map[K, V] {
	return immutable.WrapMap(m.Map())
}

func (m *Map[K, V]) SetMap(kvs c.Map[K, V]) {
	if m == nil || kvs == nil {
		return
	}
	kvs.TrackEach(func(key K, value V) { m.Set(key, value) })
}
