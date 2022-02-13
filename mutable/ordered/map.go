package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/ordered"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/slice"
)

func AsMap[K comparable, V any](elements []*map_.KV[K, V]) *Map[K, V] {
	var (
		l           = len(elements)
		uniques     = make(map[K]V, l)
		orderedKeys = make([]K, 0, l)
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		if _, ok := uniques[key]; !ok {
			orderedKeys = append(orderedKeys, key)
			uniques[key] = val
		}
	}
	return WrapMap(orderedKeys, uniques)
}

func ToMap[K comparable, V any](elements map[K]V) *Map[K, V] {
	var (
		uniques     = make(map[K]V, len(elements))
		orderedKeys = make([]K, len(elements))
	)
	for key, val := range elements {
		orderedKeys = append(orderedKeys, key)
		uniques[key] = val
	}
	return WrapMap(orderedKeys, uniques)
}

func WrapMap[K comparable, V any](orderedKeys []K, uniques map[K]V) *Map[K, V] {
	return &Map[K, V]{keys: orderedKeys, uniques: uniques, ksize: it.GetTypeSize[K]()}
}

//Map provides access to elements by key.
type Map[K comparable, V any] struct {
	keys       []K
	uniques    map[K]V
	changeMark int32
	ksize      uintptr
}

var (
	_ mutable.Settable[int, any] = (*Map[int, any])(nil)
	_ c.Map[int, any]            = (*Map[int, any])(nil)
	_ fmt.Stringer               = (*Map[int, any])(nil)
)

func (s *Map[K, V]) Begin() c.KVIterator[K, V] {
	return s.Head()
}

func (s *Map[K, V]) Head() *it.OrderedKV[K, V] {
	return it.NewOrderedKV(s.uniques, it.NewHeadS(s.keys, s.ksize))
}

func (s *Map[K, V]) Tail() *it.OrderedKV[K, V] {
	return it.NewOrderedKV(s.uniques, it.NewTailS(s.keys, s.ksize))
}

func (s *Map[K, V]) Collect() map[K]V {
	e := s.uniques
	out := make(map[K]V, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *Map[K, V]) Sort(less func(k1, k2 K) bool) *Map[K, V] {
	s.keys = slice.SortCopy(s.keys, less)
	return s
}

func (s *Map[K, V]) Len() int {
	return len(s.keys)
}

func (s *Map[K, V]) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Map[K, V]) For(walker func(*map_.KV[K, V]) error) error {
	return map_.ForOrdered(s.keys, s.uniques, walker)
}

func (s *Map[K, V]) ForEach(walker func(*map_.KV[K, V])) {
	map_.ForEachOrdered(s.keys, s.uniques, walker)
}

func (s *Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.TrackOrdered(s.keys, s.uniques, tracker)
}

func (s *Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEachOrdered(s.keys, s.uniques, tracker)
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
		e := s.keys
		u[key] = value
		s.keys = append(e, key)
		return true
	}
	return false
}

func (s *Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return s.K()
}

func (s *Map[K, V]) K() *ordered.MapKeys[K] {
	return ordered.WrapKeys(s.keys)
}

func (s *Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return s.V()
}

func (s *Map[K, V]) V() *ordered.MapValues[K, V] {
	return ordered.WrapVal(s.keys, s.uniques)
}

func (s *Map[K, V]) String() string {
	return map_.ToStringOrdered(s.keys, s.uniques)
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
