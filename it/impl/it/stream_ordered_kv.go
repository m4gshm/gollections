package it

import (
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func NewKVPipe[k comparable, v any, c any, IT typ.KVIterator[k, v]](it IT, collector collect.CollectorKV[k, v, c]) *KVIterPipe[k, v, c] {
	return &KVIterPipe[k, v, c]{it: it, collector: collector}
}

type KVIterPipe[k comparable, v any, c any] struct {
	it        typ.KVIterator[k, v]
	collector collect.CollectorKV[k, v, c]
	elements  *c
}

var _ typ.MapPipe[any, any, any] = (*KVIterPipe[any, any, any])(nil)

func (s *KVIterPipe[k, v, c]) FilterKey(fit typ.Predicate[k]) typ.MapPipe[k, v, c] {
	return NewKVPipe(FilterKV[k, v](s.it, func(key k, val v) bool { return fit(key) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) MapKey(by typ.Converter[k, k]) typ.MapPipe[k, v, c] {
	return NewKVPipe(MapKV(s.it, func(key k, val v) (k, v) { return by(key), val }), s.collector)
}

func (s *KVIterPipe[k, v, c]) FilterValue(fit typ.Predicate[v]) typ.MapPipe[k, v, c] {
	return NewKVPipe(FilterKV(s.it, func(key k, val v) bool { return fit(val) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) MapValue(by typ.Converter[v, v]) typ.MapPipe[k, v, c] {
	return NewKVPipe(MapKV(s.it, func(key k, val v) (k, v) { return key, by(val) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) Filter(fit typ.BiPredicate[k, v]) typ.MapPipe[k, v, c] {
	return NewKVPipe(FilterKV(s.it, fit), s.collector)
}

func (s *KVIterPipe[k, v, c]) Map(by typ.BiConverter[k, v, k, v]) typ.MapPipe[k, v, c] {
	return NewKVPipe(MapKV(s.it, by), s.collector)
}

func (s *KVIterPipe[k, v, c]) Track(tracker func(k, v) error) error {
	for s.it.HasNext() {
		key, val, err := s.it.Get()
		if err != nil {
			return err
		}
		if err := tracker(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *KVIterPipe[k, v, c]) TrackEach(tracker func(k, v)) error {
	return s.Track(func(key k, val v) error { tracker(key, val); return nil })
}

func (s *KVIterPipe[k, v, c]) Reduce(by op.Quaternary[k, v]) (k, v) {
	return ReduceKV(s.it, by)
}

func (s *KVIterPipe[k, v, c]) Begin() typ.KVIterator[k, v] {
	return s.it
}

func (s *KVIterPipe[k, v, c]) Collect() c {
	ref := s.elements
	var e c
	if ref == nil {
		e = s.collector(s.it)
		ref = &e
	}
	return e
}
