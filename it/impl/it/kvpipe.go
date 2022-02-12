package it

import (
	"github.com/m4gshm/gollections/c"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/op"
)

func NewKVPipe[k comparable, v any, C any, IT c.KVIterator[k, v]](it IT, collector collect.CollectorKV[k, v, C]) *KVIterPipe[k, v, C] {
	return &KVIterPipe[k, v, C]{it: it, collector: collector}
}

type KVIterPipe[k comparable, v any, C any] struct {
	it        c.KVIterator[k, v]
	collector collect.CollectorKV[k, v, C]
	elements  *C
}

var _ c.MapPipe[string, any, any] = (*KVIterPipe[string, any, any])(nil)

func (s *KVIterPipe[k, v, C]) FilterKey(fit c.Predicate[k]) c.MapPipe[k, v, C] {
	return NewKVPipe(FilterKV(s.it, func(key k, val v) bool { return fit(key) }), s.collector)
}

func (s *KVIterPipe[k, v, C]) MapKey(by c.Converter[k, k]) c.MapPipe[k, v, C] {
	return NewKVPipe(MapKV(s.it, func(key k, val v) (k, v) { return by(key), val }), s.collector)
}

func (s *KVIterPipe[k, v, C]) FilterValue(fit c.Predicate[v]) c.MapPipe[k, v, C] {
	return NewKVPipe(FilterKV(s.it, func(key k, val v) bool { return fit(val) }), s.collector)
}

func (s *KVIterPipe[k, v, C]) MapValue(by c.Converter[v, v]) c.MapPipe[k, v, C] {
	return NewKVPipe(MapKV(s.it, func(key k, val v) (k, v) { return key, by(val) }), s.collector)
}

func (s *KVIterPipe[k, v, C]) Filter(fit c.BiPredicate[k, v]) c.MapPipe[k, v, C] {
	return NewKVPipe(FilterKV(s.it, fit), s.collector)
}

func (s *KVIterPipe[k, v, C]) Map(by c.BiConverter[k, v, k, v]) c.MapPipe[k, v, C] {
	return NewKVPipe(MapKV(s.it, by), s.collector)
}

func (s *KVIterPipe[k, v, C]) Track(tracker func(k, v) error) error {
	for s.it.HasNext() {
		key, val := s.it.Get()
		if err := tracker(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *KVIterPipe[k, v, C]) Reduce(by op.Quaternary[k, v]) (k, v) {
	return ReduceKV(s.it, by)
}

func (s *KVIterPipe[k, v, C]) Begin() c.KVIterator[k, v] {
	return s.it
}

func (s *KVIterPipe[k, v, C]) Collect() C {
	ref := s.elements
	var e C
	if ref == nil {
		e = s.collector(s.it)
		ref = &e
	}
	return e
}
