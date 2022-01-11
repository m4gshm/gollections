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

var _ mutable.Map[any, any, typ.Iterator[*typ.KV[any, any]]] = (*OrderedMap[any, any])(nil)
var _ typ.Map[any, any, typ.Iterator[*typ.KV[any, any]]] = (*OrderedMap[any, any])(nil)

// var _ fmt.Stringer = (*OrderedMap[interface{}, interface{}])(nil)

func (s *OrderedMap[k, v]) Begin() typ.Iterator[*typ.KV[k, v]] {
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

func (s *OrderedMap[k, v]) TrackEach(tracker func(k, v)) error {
	e := s.uniques
	for _, ref := range s.elements {
		key := *ref
		tracker(key, e[key])
	}
	return nil
}

func (s *OrderedMap[k, v]) ForEach(walker func(*typ.KV[k, v])) error {
	return s.TrackEach(func(k k, v v) { walker(K.V(k, v)) })
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

func (s *OrderedMap[k, v]) Keys() typ.Container[k, []k, int, typ.Iterator[k]] {
	return ref.Wrap(s.elements)
}

func (s *OrderedMap[k, v]) Values() typ.Container[v, []v, int, typ.Iterator[v]] {
	return dict.Wrap(s.elements, s.uniques)
}

func (s *OrderedMap[k, v]) Filter(filter typ.Predicate[*typ.KV[k, v]]) typ.Pipe[*typ.KV[k, v], map[k]v, typ.Iterator[*typ.KV[k, v]]] {
	return it.NewOrderedKVPipe[k, v](it.Filter(s.Iter(), filter), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) Map(by typ.Converter[*typ.KV[k, v], *typ.KV[k, v]]) typ.Pipe[*typ.KV[k, v], map[k]v, typ.Iterator[*typ.KV[k, v]]] {
	return it.NewOrderedKVPipe[k, v](it.Map(s.Iter(), by), collect.Map[k, v])
}

func (s *OrderedMap[k, v]) Reduce(by op.Binary[*typ.KV[k, v]]) *typ.KV[k, v] {
	return it.Reduce(s.Iter(), by)
}
