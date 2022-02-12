package ordered

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/slice"

	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
)

func ConvertKVsToMap[K comparable, V any](elements []*map_.KV[K, V]) *Map[K, V] {
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

func WrapMap[K comparable, V any](order []K, uniques map[K]V) *Map[K, V] {
	return &Map[K, V]{elements: order, uniques: uniques}
}

type Map[K comparable, V any] struct {
	elements []K
	uniques  map[K]V
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ fmt.Stringer    = (*Map[int, any])(nil)
)

func (s *Map[K, V]) Begin() c.KVIterator[K, V] {
	return s.Head()
}

func (s *Map[K, V]) Head() *it.OrderedKV[K, V] {
	return it.NewOrderedKV(s.elements, s.uniques, it.NewHead[K])
}

func (s *Map[K, V]) Tail() *it.OrderedKV[K, V] {
	return it.NewOrderedKV(s.elements, s.uniques, it.NewTail[K])
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
	return WrapMap(slice.SortCopy(s.elements, less), s.uniques)
}

func (s *Map[K, V]) String() string {
	return map_.ToStringOrdered(s.elements, s.uniques)
}

func (s *Map[K, V]) Track(tracker func(K, V) error) error {
	return map_.TrackOrdered(s.elements, s.uniques, tracker)
}

func (s *Map[K, V]) TrackEach(tracker func(K, V)) {
	map_.TrackEachOrdered(s.elements, s.uniques, tracker)
}

func (s *Map[K, V]) For(walker func(*map_.KV[K, V]) error) error {
	return map_.ForOrdered(s.elements, s.uniques, walker)
}

func (s *Map[K, V]) ForEach(walker func(*map_.KV[K, V])) {
	map_.ForEachOrdered(s.elements, s.uniques, walker)
}

func (s *Map[K, V]) Len() int {
	return len(s.elements)
}

func (s *Map[K, V]) Contains(key K) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *Map[K, V]) Get(key K) (V, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *Map[K, V]) Keys() c.Collection[K, []K, c.Iterator[K]] {
	return s.K()
}

func (s *Map[K, V]) K() *MapKeys[K] {
	return WrapKeys(s.elements)
}

func (s *Map[K, V]) Values() c.Collection[V, []V, c.Iterator[V]] {
	return s.V()
}

func (s *Map[K, V]) V() *MapValues[K, V] {
	return WrapVal(s.elements, s.uniques)
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
