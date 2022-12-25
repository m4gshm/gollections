package it

import (
	"github.com/m4gshm/gollections/c"
)

// NewKVPipe instantiates Iterator wrapper that converts the elements into key/value pairs and iterates over them.
func NewKVPipe[K comparable, V any, C any, IT c.KVIterator[K, V]](it IT, collector KVCollector[K, V, C]) *KVIterPipe[K, V, C] {
	return &KVIterPipe[K, V, C]{it: it, collector: collector}
}

// KVIterPipe is the key/value Iterator based pipe implementation.
type KVIterPipe[K comparable, V any, C any] struct {
	it        c.KVIterator[K, V]
	collector KVCollector[K, V, C]
	out       *C
}

var _ c.MapPipe[string, any, any] = (*KVIterPipe[string, any, any])(nil)

func (s *KVIterPipe[K, V, C]) FilterKey(fit c.Predicate[K]) c.MapPipe[K, V, C] {
	var kvFit c.BiPredicate[K, V] = func(key K, val V) bool { return fit(key) }
	return NewKVPipe(FilterKV(s.it, kvFit), s.collector)
}

func (s *KVIterPipe[K, V, C]) MapKey(by c.Converter[K, K]) c.MapPipe[K, V, C] {
	return NewKVPipe(MapKV(s.it, func(key K, val V) (K, V) { return by(key), val }), s.collector)
}

func (s *KVIterPipe[K, V, C]) FilterValue(fit c.Predicate[V]) c.MapPipe[K, V, C] {
	var kvFit c.BiPredicate[K, V] = func(key K, val V) bool { return fit(val) }
	return NewKVPipe(FilterKV(s.it, kvFit), s.collector)
}

func (s *KVIterPipe[K, V, C]) MapValue(by c.Converter[V, V]) c.MapPipe[K, V, C] {
	return NewKVPipe(MapKV(s.it, func(key K, val V) (K, V) { return key, by(val) }), s.collector)
}

func (s *KVIterPipe[K, V, C]) Filter(fit c.BiPredicate[K, V]) c.MapPipe[K, V, C] {
	return NewKVPipe(FilterKV(s.it, fit), s.collector)
}

func (s *KVIterPipe[K, V, C]) Map(by c.BiConverter[K, V, K, V]) c.MapPipe[K, V, C] {
	return NewKVPipe(MapKV(s.it, by), s.collector)
}

func (s *KVIterPipe[K, V, C]) Track(tracker func(K, V) error) error {
	for key, val, ok := s.it.Next(); ok; key, val, ok = s.it.Next() {
		if err := tracker(key, val); err != nil {
			return err
		}
	}
	return nil
}

func (s *KVIterPipe[K, V, C]) Reduce(by c.Quaternary[K, V]) (K, V) {
	return ReduceKV(s.it, by)
}

func (s *KVIterPipe[K, V, C]) Begin() c.KVIterator[K, V] {
	return s.it
}

func (s *KVIterPipe[K, V, C]) Collect() C {
	var e C
	if s.out == nil {
		e = s.collector(s.it)
		s.out = &e
	}
	return e
}

// KVCollector is Converter of key/value Iterator that collects all values to any slice or map, mostly used to extract slice fields to flatting a result
type KVCollector[k, v any, out any] func(c.KVIterator[k, v]) out
