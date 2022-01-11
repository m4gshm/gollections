package it

import (
	"github.com/m4gshm/gollections/collect"
	"github.com/m4gshm/gollections/op"
	"github.com/m4gshm/gollections/typ"
)

func NewOrderedKVPipe[k comparable, v any, c any, IT typ.Iterator[*typ.KV[k, v]]](it IT, collector collect.Collector[*typ.KV[k, v], c]) *KVOrdererIterPipe[k, v, c] {
	return &KVOrdererIterPipe[k, v, c]{it: it, collector: collector}
}

type KVOrdererIterPipe[k comparable, v any, c any] struct {
	it        typ.Iterator[*typ.KV[k, v]]
	collector collect.Collector[*typ.KV[k, v], c]
	elements  *c
}

var _ typ.Pipe[*typ.KV[any, any], any, typ.Iterator[*typ.KV[any, any]]] = (*KVOrdererIterPipe[any, any, any])(nil)

func (s *KVOrdererIterPipe[k, v, cnt]) Filter(fit typ.Predicate[*typ.KV[k, v]]) typ.Pipe[*typ.KV[k, v], cnt, typ.Iterator[*typ.KV[k, v]]] {
	return NewOrderedKVPipe(Filter(s.it, fit), s.collector)
}

func (s *KVOrdererIterPipe[k, v, c]) Map(by typ.Converter[*typ.KV[k, v], *typ.KV[k, v]]) typ.Pipe[*typ.KV[k, v], c, typ.Iterator[*typ.KV[k, v]]] {
	return NewOrderedKVPipe(Map(s.it, by), s.collector)
}

func (s *KVOrdererIterPipe[k, v, c]) TrackEach(tracker func(k, v)) error {
	return s.ForEach(func(kv *typ.KV[k, v]) { tracker(kv.Key(), kv.Value()) })
}

func (s *KVOrdererIterPipe[k, v, c]) ForEach(walker func(*typ.KV[k, v])) error {
	for s.it.HasNext() {
		n, err := s.it.Next()
		if err != nil {
			return err
		}
		walker(n)
	}
	return nil
}

func (s *KVOrdererIterPipe[k, v, c]) Reduce(by op.Binary[*typ.KV[k, v]]) *typ.KV[k, v] {
	return Reduce(s.it, by)
}

func (s *KVOrdererIterPipe[k, v, c]) Begin() typ.Iterator[*typ.KV[k, v]] {
	return s.it
}

func (s *KVOrdererIterPipe[k, v, c]) Collect() c {
	ref := s.elements
	var e c
	if ref == nil {
		e = s.collector(s.it)
		ref = &e
	}
	return e
}
