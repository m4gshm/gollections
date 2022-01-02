package dict

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

func NewOrderedMap[k comparable, v any](elements []*typ.KV[k, v]) *OrderedMap[k, v] {
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

type OrderedMap[k comparable, v any] struct {
	elements []*k
	uniques  map[k]v
}

var _ typ.Map[any, any] = (*OrderedMap[interface{}, interface{}])(nil)

// var _ fmt.Stringer = (*OrderedMap[interface{}, interface{}])(nil)

func (s *OrderedMap[k, v]) Begin() typ.Iterator[*typ.KV[k, v]] {
	return iter.NewOrderKV(s.elements, s.uniques)
}

func (s *OrderedMap[k, v]) Values() map[k]v {
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
