package omap

import (
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/vector/ref"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func Convert[k comparable, v any](elements []*typ.KV[k, v]) *OrderedMap[k, v] {
	var (
		uniques = make(map[k]v, 0)
		order   = make([]*k, 0, 0)
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		if _, ok := uniques[key]; !ok {
			order = append(order, &key)
			uniques[key] = val
		}
	}
	return Wrap(order, uniques)
}

func ConvertMap[k comparable, v any](elements map[k]v) *OrderedMap[k, v] {
	var (
		uniques = make(map[k]v, len(elements))
		order   = make([]*k, len(elements))
	)
	for key, val := range elements {
		order = append(order, &key)
		uniques[key] = val
	}
	return Wrap(order, uniques)
}

func Wrap[k comparable, v any](order []*k, uniques map[k]v) *OrderedMap[k, v] {
	return &OrderedMap[k, v]{elements: order, uniques: uniques}
}

type OrderedMap[k comparable, v any] struct {
	elements []*k
	uniques  map[k]v
}

var _ typ.Map[any, any, typ.KVIterator[any, any]] = (*OrderedMap[any, any])(nil)

// var _ fmt.Stringer = (*OrderedMap[interface{}, interface{}])(nil)

func (s *OrderedMap[k, v]) Begin() typ.KVIterator[k, v] {
	return s.Iter()
}

func (s *OrderedMap[k, v]) Iter() *it.OrderedKV[k, v] {
	return it.NewOrderedKV(s.elements, s.uniques)
}

func (s *OrderedMap[k, v]) Collect() map[k]v {
	e := s.uniques
	out := make(map[k]v, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *OrderedMap[k, v]) Track(tracker func(k, v) error) error {
	return map_.TrackOrdered(s.elements, s.uniques, tracker)
}

func (s *OrderedMap[k, v]) TrackEach(tracker func(k, v)) error {
	return map_.TrackEachOrdered(s.elements, s.uniques, tracker)
}

func (s *OrderedMap[k, v]) For(walker func(*typ.KV[k, v]) error) error {
	return map_.ForOrdered(s.elements, s.uniques, walker)
}

func (s *OrderedMap[k, v]) ForEach(walker func(*typ.KV[k, v])) error {
	return map_.ForEachOrdered(s.elements, s.uniques, walker)
}

func (s *OrderedMap[k, v]) Len() int {
	return len(s.elements)
}

func (s *OrderedMap[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *OrderedMap[k, v]) Get(key k) (v, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *OrderedMap[k, v]) Keys() typ.Collection[k, []k, typ.Iterator[k]] {
	return ref.Wrap(s.elements)
}

func (s *OrderedMap[k, v]) Values() typ.Collection[v, []v, typ.Iterator[v]] {
	return WrapVal(s.elements, s.uniques)
}

func (s *OrderedMap[k, v]) FilterKey(fit typ.Predicate[k]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(key) }), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) MapKey(by typ.Converter[k, k]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return by(key), val }), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) FilterValue(fit typ.Predicate[v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(val) }), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) MapValue(by typ.Converter[v, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return key, by(val) }), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) Filter(filter typ.BiPredicate[k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), filter), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) Map(by typ.BiConverter[k, v, k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), by), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) Reduce(by op.Quaternary[k, v]) (k, v) {
	return it.ReduceKV(s.Iter(), by)
}
