package map_

import (
	"fmt"

	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/set"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func Convert[k comparable, v any](elements []*typ.KV[k, v]) *Map[k, v] {
	uniques := make(map[k]v, 0)
	for _, kv := range elements {
		uniques[kv.Key()] = kv.Value()
	}
	return Wrap(uniques)
}

func ConvertMap[k comparable, v any](elements map[k]v) *Map[k, v] {
	uniques := make(map[k]v, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return Wrap(uniques)
}

func Wrap[k comparable, v any](uniques map[k]v) *Map[k, v] {
	return &Map[k, v]{uniques: uniques}
}

type Map[k comparable, v any] struct {
	uniques map[k]v
}

var _ typ.Map[any, any, typ.KVIterator[any, any]] = (*Map[any, any])(nil)

var _ fmt.Stringer = (*Map[any, any])(nil)

func (s *Map[k, v]) Begin() typ.KVIterator[k, v] {
	return s.Iter()
}

func (s *Map[k, v]) Iter() *it.KV[k, v] {
	return it.NewKV(s.uniques)
}

func (s *Map[k, v]) Collect() map[k]v {
	e := s.uniques
	out := make(map[k]v, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *Map[k, v]) Track(tracker func(k, v) error) error {
	return map_.Track(s.uniques, tracker)
}

func (s *Map[k, v]) TrackEach(tracker func(k, v)) {
	map_.TrackEach(s.uniques, tracker)
}

func (s *Map[k, v]) For(walker func(*typ.KV[k, v]) error) error {
	return map_.For(s.uniques, walker)
}

func (s *Map[k, v]) ForEach(walker func(*typ.KV[k, v])) {
	map_.ForEach(s.uniques, walker)
}

func (s *Map[k, v]) Len() int {
	return len(s.uniques)
}

func (s *Map[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *Map[k, v]) Get(key k) (v, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *Map[k, v]) Keys() typ.Collection[k, []k, typ.Iterator[k]] {
	return set.Wrap(s.uniques)
}

func (s *Map[k, v]) Values() typ.Collection[v, []v, typ.Iterator[v]] {
	return WrapVal(s.uniques)
}

func (s *Map[k, v]) String() string {
	return map_.ToString(s.uniques)
}

func (s *Map[k, v]) FilterKey(fit typ.Predicate[k]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(key) }), collect.Map[k, v])
}

func (s *Map[k, v]) MapKey(by typ.Converter[k, k]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return by(key), val }), collect.Map[k, v])
}

func (s *Map[k, v]) FilterValue(fit typ.Predicate[v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(val) }), collect.Map[k, v])
}

func (s *Map[k, v]) MapValue(by typ.Converter[v, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return key, by(val) }), collect.Map[k, v])
}

func (s *Map[k, v]) Filter(filter typ.BiPredicate[k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), filter), collect.Map[k, v])
}

func (s *Map[k, v]) Map(by typ.BiConverter[k, v, k, v]) typ.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), by), collect.Map[k, v])
}

func (s *Map[k, v]) Reduce(by op.Quaternary[k, v]) (k, v) {
	return it.ReduceKV(s.Iter(), by)
}
