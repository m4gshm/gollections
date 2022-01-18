package map_

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/vector/dict"
	"github.com/m4gshm/gollections/immutable/vector/ref"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func ToOrderedMap[k comparable, v any](elements []*typ.KV[k, v]) *OrderedMap[k, v] {
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
	return &OrderedMap[k, v]{elements: order, uniques: uniques}
}

func NewOrderedMap[k comparable, v any](capacity int) *OrderedMap[k, v] {
	return &OrderedMap[k, v]{elements: make([]*k, 0, capacity), uniques: make(map[k]v, capacity)}
}

type OrderedMap[k comparable, v any] struct {
	changeMark int32
	elements   []*k
	uniques    map[k]v
	err        error
}

var _ mutable.Map[any, any, typ.KVIterator[any, any]] = (*OrderedMap[any, any])(nil)
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

func (s *OrderedMap[k, v]) For(walker func(*typ.KV[k, v]) error) error {
	return s.Track(func(key k, value v) error { return walker(K.V(key, value)) })
}

func (s *OrderedMap[k, v]) ForEach(walker func(*typ.KV[k, v])) error {
	return s.For(func(kv *typ.KV[k, v]) error { walker(kv); return nil })
}

func (s *OrderedMap[k, v]) Track(tracker func(k, v) error) error {
	e := s.uniques
	for _, ref := range s.elements {
		key := *ref
		if err := tracker(key, e[key]); err != nil {
			return err
		}
	}
	return nil
}

func (s *OrderedMap[k, v]) TrackEach(tracker func(k, v)) error {
	return s.Track(func(key k, value v) error { tracker(key, value); return nil })
}

func (s *OrderedMap[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *OrderedMap[k, v]) Get(key k) (v, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *OrderedMap[k, v]) Set(key k, value v) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	u := s.uniques
	if _, ok := u[key]; !ok {
		markOnStart := s.changeMark
		e := s.elements
		u[key] = value
		s.elements = append(e, &key)
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *OrderedMap[k, v]) Keys() typ.Container[[]k, typ.Iterator[k]] {
	return ref.Wrap(s.elements)
}

func (s *OrderedMap[k, v]) Values() typ.Container[[]v, typ.Iterator[v]] {
	return dict.Wrap(s.elements, s.uniques)
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
