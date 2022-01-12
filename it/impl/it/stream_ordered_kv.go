package it

import (
	"github.com/m4gshm/gollections/K"
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func NewKVPipe[k comparable, v any, c any, IT typ.Iterator[*typ.KV[k, v]]](it IT, collector collect.Collector[*typ.KV[k, v], c]) *KVIterPipe[k, v, c] {
	return &KVIterPipe[k, v, c]{it: it, collector: collector}
}

type KVIterPipe[k comparable, v any, c any] struct {
	it        typ.Iterator[*typ.KV[k, v]]
	collector collect.Collector[*typ.KV[k, v], c]
	elements  *c
}

// var _ typ.Pipe[*typ.KV[any, any], any, typ.Iterator[*typ.KV[any, any]]] = (*KVIterPipe[any, any, any])(nil)
var _ typ.MapPipe[any, any, any] = (*KVIterPipe[any, any, any])(nil)

func (s *KVIterPipe[k, v, c]) FilterKey(fit typ.Predicate[k]) typ.MapPipe[k, v, c] {
	return NewKVPipe(Filter(s.it, func(kv *typ.KV[k, v]) bool { return fit(kv.Key()) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) MapKey(by typ.Converter[k, k]) typ.MapPipe[k, v, c] {
	return NewKVPipe(Map(s.it, func(kv *typ.KV[k, v]) *typ.KV[k, v] { return K.V(by(kv.Key()), kv.Value()) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) FilterValue(fit typ.Predicate[v]) typ.MapPipe[k, v, c] {
	return NewKVPipe(Filter(s.it, func(kv *typ.KV[k, v]) bool { return fit(kv.Value()) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) MapValue(by typ.Converter[v, v]) typ.MapPipe[k, v, c] {
	return NewKVPipe(Map(s.it, func(kv *typ.KV[k, v]) *typ.KV[k, v] { return K.V(kv.Key(), by(kv.Value())) }), s.collector)
}

func (s *KVIterPipe[k, v, c]) Filter(fit typ.Predicate[*typ.KV[k, v]]) typ.MapPipe[k, v, c] {
	return NewKVPipe(Filter(s.it, fit), s.collector)
}

func (s *KVIterPipe[k, v, c]) Map(by typ.Converter[*typ.KV[k, v], *typ.KV[k, v]]) typ.MapPipe[k, v, c] {
	return NewKVPipe(Map(s.it, by), s.collector)
}

func (s *KVIterPipe[k, v, c]) TrackEach(tracker func(k, v)) error {
	return s.ForEach(func(kv *typ.KV[k, v]) { tracker(kv.Key(), kv.Value()) })
}

func (s *KVIterPipe[k, v, c]) ForEach(walker func(*typ.KV[k, v])) error {
	for s.it.HasNext() {
		n, err := s.it.Next()
		if err != nil {
			return err
		}
		walker(n)
	}
	return nil
}

func (s *KVIterPipe[k, v, c]) Reduce(by op.Binary[*typ.KV[k, v]]) *typ.KV[k, v] {
	return Reduce(s.it, by)
}

func (s *KVIterPipe[k, v, c]) Begin() typ.Iterator[*typ.KV[k, v]] {
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
