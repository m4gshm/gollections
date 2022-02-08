package immutable

import (
	"fmt"

	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/map_"
	m "github.com/m4gshm/gollections/map_"
	"github.com/m4gshm/gollections/op"
)

func ConvertKVsToMap[k comparable, v any](elements []*map_.KV[k, v]) *Map[k, v] {
	uniques := make(map[k]v, len(elements))
	for _, kv := range elements {
		uniques[kv.Key()] = kv.Value()
	}
	return WrapMap(uniques)
}

func NewMap[k comparable, v any](elements map[k]v) *Map[k, v] {
	uniques := make(map[k]v, len(elements))
	for key, val := range elements {
		uniques[key] = val
	}
	return WrapMap(uniques)
}

func WrapMap[k comparable, v any](uniques map[k]v) *Map[k, v] {
	return &Map[k, v]{uniques: uniques}
}

//Map provides access to elements by key.
type Map[k comparable, v any] struct {
	uniques map[k]v
}

var (
	_ c.Map[int, any] = (*Map[int, any])(nil)
	_ fmt.Stringer    = (*Map[int, any])(nil)
)

func (s *Map[k, v]) Begin() c.KVIterator[k, v] {
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

func (s *Map[k, v]) String() string {
	return m.ToString(s.uniques)
}

func (s *Map[k, v]) Track(tracker func(k, v) error) error {
	return m.Track(s.uniques, tracker)
}

func (s *Map[k, v]) TrackEach(tracker func(k, v)) {
	m.TrackEach(s.uniques, tracker)
}

func (s *Map[k, v]) For(walker func(*map_.KV[k, v]) error) error {
	return m.For(s.uniques, walker)
}

func (s *Map[k, v]) ForEach(walker func(*map_.KV[k, v])) {
	m.ForEach(s.uniques, walker)
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

func (s *Map[k, v]) Keys() c.Collection[k, []k, c.Iterator[k]] {
	return WrapKeys(s.uniques)
}

func (s *Map[k, v]) Values() c.Collection[v, []v, c.Iterator[v]] {
	return WrapVal(s.uniques)
}

func (s *Map[k, v]) FilterKey(fit c.Predicate[k]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(key) }), collect.Map[k, v])
}

func (s *Map[k, v]) MapKey(by c.Converter[k, k]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return by(key), val }), collect.Map[k, v])
}

func (s *Map[k, v]) FilterValue(fit c.Predicate[v]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), func(key k, val v) bool { return fit(val) }), collect.Map[k, v])
}

func (s *Map[k, v]) MapValue(by c.Converter[v, v]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), func(key k, val v) (k, v) { return key, by(val) }), collect.Map[k, v])
}

func (s *Map[k, v]) Filter(filter c.BiPredicate[k, v]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.FilterKV(s.Iter(), filter), collect.Map[k, v])
}

func (s *Map[k, v]) Map(by c.BiConverter[k, v, k, v]) c.MapPipe[k, v, map[k]v] {
	return it.NewKVPipe(it.MapKV(s.Iter(), by), collect.Map[k, v])
}

func (s *Map[k, v]) Reduce(by op.Quaternary[k, v]) (k, v) {
	return it.ReduceKV(s.Iter(), by)
}
