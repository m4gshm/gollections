package mutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	m "github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
)

//NewMap creates the Map with a predefined capacity.
func NewMap[K comparable, V any](capacity int) Map[K, V] {
	return WrapMap(make(map[K]V, capacity))
}

//AsMap converts a slice of key/value pairs to teh Map.
func AsMap[K comparable, V any](elements []c.KV[K, V]) Map[K, V] {
	var (
		uniques = make(map[K]V, len(elements))
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		uniques[key] = val
	}
	return WrapMap(uniques)
}

//ToMap creates the Map and copies elements to it.
func ToMap[K comparable, V any](elements map[K]V) Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

//WrapMap creates the Map using a map as internal storage.
func WrapMap[K comparable, V any](elements map[K]V) Map[K, V] {
	return Map[K, V](elements)
}

//Map is the Collection implementation based on the embedded map.
type Map[K comparable, V any] map[K]V

var (
	_ c.Settable[int, any] = (*Map[int, any])(nil)
	_ c.Settable[int, any] = (Map[int, any])(nil)
	_ c.Map[int, any]      = (*Map[int, any])(nil)
	_ c.Map[int, any]      = (Map[int, any])(nil)
	_ fmt.Stringer         = (*Map[int, any])(nil)
	_ fmt.Stringer         = (Map[int, any])(nil)
)

func (s Map[K, V]) Begin() c.KVIterator[K, V] {
	iter := s.Head()
	return &iter
}

func (s Map[K, V]) Head() it.KV[K, V] {
	return it.NewKV(s)
}

func (s Map[K, V]) First() (it.KV[K, V], K, V, bool) {
	var (
		iter               = it.NewKV(s)
		firstK, firstV, ok = iter.Next()
	)
	return iter, firstK, firstV, ok
}

func (s Map[K, V]) Collect() map[K]V {
	return map_.Copy(s)
}

func (s Map[K, V]) Len() int {
	return len(s)
}

func (s Map[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s Map[K, V]) For(walker func(c.KV[K, V]) error) error {
	return m.For(s, walker)
}

func (s Map[K, V]) ForEach(walker func(c.KV[K, V])) {
	m.ForEach(s, walker)
}

func (s Map[K, V]) Track(tracker func(K, V) error) error {
	return m.Track(s, tracker)
}

func (s Map[K, V]) TrackEach(tracker func(K, V)) {
	m.TrackEach(s, tracker)
}

func (s Map[K, V]) Contains(key K) bool {
	_, ok := s[key]
	return ok
}

func (s Map[K, V]) Get(key K) (V, bool) {
	val, ok := s[key]
	return val, ok
}

func (s Map[K, V]) Set(key K, value V) bool {
	u := s
	if _, ok := u[key]; !ok {
		u[key] = value
		return true
	}
	return false
}

func (s Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return s.K()
}

func (s Map[K, V]) K() immutable.MapKeys[K, V] {
	return immutable.WrapKeys(s)
}

func (s Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return s.V()
}

func (s Map[K, V]) V() immutable.MapValues[K, V] {
	return immutable.WrapVal(s)
}

func (s Map[K, V]) String() string {
	return m.ToString(s)
}

func (s Map[K, V]) FilterKey(fit c.Predicate[K]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.FilterKV(&iter, func(key K, val V) bool { return fit(key) }), collect.Map[K, V])
}

func (s Map[K, V]) MapKey(by c.Converter[K, K]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.MapKV(&iter, func(key K, val V) (K, V) { return by(key), val }), collect.Map[K, V])
}

func (s Map[K, V]) FilterValue(fit c.Predicate[V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.FilterKV(&iter, func(key K, val V) bool { return fit(val) }), collect.Map[K, V])
}

func (s Map[K, V]) MapValue(by c.Converter[V, V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.MapKV(&iter, func(key K, val V) (K, V) { return key, by(val) }), collect.Map[K, V])
}

func (s Map[K, V]) Filter(filter c.BiPredicate[K, V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.FilterKV(&iter, filter), collect.Map[K, V])
}

func (s Map[K, V]) Map(by c.BiConverter[K, V, K, V]) c.MapPipe[K, V, map[K]V] {
	iter := s.Head()
	return it.NewKVPipe(it.MapKV(&iter, by), collect.Map[K, V])
}

func (s Map[K, V]) Reduce(by op.Quaternary[K, V]) (K, V) {
	iter := s.Head()
	return it.ReduceKV(&iter, by)
}
