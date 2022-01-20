package omap

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/immutable/map_"
	"github.com/m4gshm/gollections/it/impl/it"
	"github.com/m4gshm/gollections/mutable"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func Convert[k comparable, v any](elements []*typ.KV[k, v]) *Map[k, v] {
	var (
		uniques = make(map[k]v, 0)
	)
	for _, kv := range elements {
		key := kv.Key()
		val := kv.Value()
		uniques[key] = val
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
	uniques    map[k]v
	err        error
	changeMark int32
}

var _ mutable.Map[any, any, typ.KVIterator[any, any]] = (*Map[any, any])(nil)
var _ typ.Map[any, any, typ.KVIterator[any, any]] = (*Map[any, any])(nil)

// var _ fmt.Stringer = (*Map[interface{}, interface{}])(nil)

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

func (s *Map[k, v]) Len() int {
	return len(s.uniques)
}

func (s *Map[k, v]) For(walker func(*typ.KV[k, v]) error) error {
	return s.Track(func(key k, value v) error { return walker(K.V(key, value)) })
}

func (s *Map[k, v]) ForEach(walker func(*typ.KV[k, v])) error {
	return s.For(func(kv *typ.KV[k, v]) error { walker(kv); return nil })
}

func (s *Map[k, v]) Track(tracker func(k, v) error) error {
	uniques := s.uniques
	for key, val := range uniques {
		if err := tracker(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *Map[k, v]) TrackEach(tracker func(k, v)) error {
	return s.Track(func(key k, value v) error { tracker(key, value); return nil })
}

func (s *Map[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *Map[k, v]) Get(key k) (v, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *Map[k, v]) Set(key k, value v) (bool, error) {
	if err := s.err; err != nil {
		return false, err
	}
	u := s.uniques
	if _, ok := u[key]; !ok {
		markOnStart := s.changeMark
		u[key] = value
		return mutable.Commit(markOnStart, &s.changeMark, &s.err)
	}
	return false, nil
}

func (s *Map[k, v]) Keys() typ.Collection[k, []k, typ.Iterator[k]] {
	return map_.WrapKeys(s.uniques)
}

func (s *Map[k, v]) Values() typ.Collection[v, []v, typ.Iterator[v]] {
	return map_.WrapVal(s.uniques)
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
