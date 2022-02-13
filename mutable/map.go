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

func AsMap[K comparable, V any](elements []*map_.KV[K, V]) *Map[K, V] {
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

func ToMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	uniques := make(map[K]V, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

func WrapMap[K comparable, V any](uniques map[K]V) *Map[K, V] {
	return &Map[K, V]{uniques: uniques}
}

type Map[K comparable, V any] struct {
	uniques map[K]V
}

var _ Settable[int, any] = (*Map[int, any])(nil)
var _ c.Map[int, any] = (*Map[int, any])(nil)
var _ fmt.Stringer = (*Map[int, any])(nil)

func (s *Map[K, V]) Begin() c.KVIterator[K, V] {
	return s.Head()
}

func (s *Map[K, V]) Head() *it.KV[K, V] {
	return it.NewKV(s.uniques)
}

func (s *Map[K, V]) Collect() map[K]V {
	e := s.uniques
	out := make(map[K]V, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *Map[K, V]) Len() int {
	return len(s.uniques)
}

func (s *Map[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Map[K, V]) For(walker func(*map_.KV[K, V]) error) error {
	return m.For(s.uniques, walker)
}

func (s *Map[K, V]) ForEach(walker func(*map_.KV[K, V])) {
	m.ForEach(s.uniques, walker)
}

func (s *Map[K, V]) Track(tracker func(K, V) error) error {
	return m.Track(s.uniques, tracker)
}

func (s *Map[K, V]) TrackEach(tracker func(K, V)) {
	m.TrackEach(s.uniques, tracker)
}

func (s *Map[K, V]) Contains(key K) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *Map[K, V]) Get(key K) (V, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *Map[K, V]) Set(key K, value V) bool {
	u := s.uniques
	if _, ok := u[key]; !ok {
		u[key] = value
		return true
	}
	return false
}

func (s *Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return s.K()
}

func (s *Map[K, V]) K() *immutable.MapKeys[K, V] {
	return immutable.WrapKeys(s.uniques)
}

func (s *Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return s.V()
}

func (s *Map[K, V]) V() *immutable.MapValues[K, V] {
	return immutable.WrapVal(s.uniques)
}

func (s *Map[K, V]) String() string {
	return m.ToString(s.uniques)
}

func (s *Map[K, V]) FilterKey(fit c.Predicate[K]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(s.Head(), func(key K, val V) bool { return fit(key) }), collect.Map[K, V])
}

func (s *Map[K, V]) MapKey(by c.Converter[K, K]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.MapKV(s.Head(), func(key K, val V) (K, V) { return by(key), val }), collect.Map[K, V])
}

func (s *Map[K, V]) FilterValue(fit c.Predicate[V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(s.Head(), func(key K, val V) bool { return fit(val) }), collect.Map[K, V])
}

func (s *Map[K, V]) MapValue(by c.Converter[V, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.MapKV(s.Head(), func(key K, val V) (K, V) { return key, by(val) }), collect.Map[K, V])
}

func (s *Map[K, V]) Filter(filter c.BiPredicate[K, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.FilterKV(s.Head(), filter), collect.Map[K, V])
}

func (s *Map[K, V]) Map(by c.BiConverter[K, V, K, V]) c.MapPipe[K, V, map[K]V] {
	return it.NewKVPipe(it.MapKV(s.Head(), by), collect.Map[K, V])
}

func (s *Map[K, V]) Reduce(by op.Quaternary[K, V]) (K, V) {
	return it.ReduceKV(s.Head(), by)
}
