package dict

import (
	"github.com/m4gshm/container/immutable/vector/dict"
	"github.com/m4gshm/container/immutable/vector/ref"
	"github.com/m4gshm/container/it/impl/it"
	"github.com/m4gshm/container/mutable"
	"github.com/m4gshm/container/typ"
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
	elements []*k
	uniques  map[k]v
}

var _ mutable.Map[any, any, typ.Iterator[*typ.KV[any, any]]] = (*OrderedMap[any, any])(nil)

// var _ fmt.Stringer = (*OrderedMap[interface{}, interface{}])(nil)

func (s *OrderedMap[k, v]) Begin() typ.Iterator[*typ.KV[k, v]] {
	return s.Iter()
}

func (s *OrderedMap[k, v]) Iter() *it.OrderedKV[k, v] {
	return it.NewOrderedKV(s.elements, s.uniques)
}

func (s *OrderedMap[k, v]) Elements() map[k]v {
	e := s.uniques
	out := make(map[k]v, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *OrderedMap[k, v]) ForEach(tracker func(k, v)) {
	e := s.uniques
	for _, ref := range s.elements {
		key := *ref
		tracker(key, e[key])
	}
}

func (s *OrderedMap[k, v]) Len() int {
	return len(s.elements)
}

func (s *OrderedMap[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}

func (s *OrderedMap[k, v]) Get(key k) (v, bool) {
	val, ok := s.uniques[key]
	return val, ok
}

func (s *OrderedMap[k, v]) Put(key k, value v) bool {
	u := s.uniques
	if _, ok := u[key]; !ok {
		e := s.elements
		u[key] = value
		s.elements = append(e, &key)
		return true
	}
	return false
}

func (s *OrderedMap[k, v]) Keys() typ.Container[k, int, typ.Iterator[k]] {
	return ref.Wrap(s.elements)
}

func (s *OrderedMap[k, v]) Values() typ.Container[v, int, typ.Iterator[v]] {
	return dict.Wrap(s.elements, s.uniques)
}
