package immutable

import (
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

func newMap[k comparable, v any](values []*typ.KV[k, v]) *OrderedMap[k, v] {
	var (
		uniques = make(map[k]v, 0)
		order   = make([]*k, 0, 0)
	)
	for _, kv := range values {
		key := kv.Key()
		val := kv.Value()
		if _, ok := uniques[key]; !ok {
			order = append(order, &key)
			uniques[key] = val
		}
	}
	return &OrderedMap[k, v]{order: order, uniques: uniques}
}

type OrderedMap[k comparable, v any] struct {
	order   []*k
	uniques map[k]v
}

var _ Map[any, any, *iter.OrderedKVIter[any, any]] = (*OrderedMap[interface{}, interface{}])(nil)
// var _ fmt.Stringer = (*OrderedMap[interface{}, interface{}])(nil)

func (s *OrderedMap[k, v]) Begin() *iter.OrderedKVIter[k, v] {
	return iter.NewOrderKV(&s.order, s.uniques)
}

func (s *OrderedMap[k, v]) Values() map[k]v {
	e := s.uniques
	out := make(map[k]v, len(e))
	for key, val := range e {
		out[key] = val
	}
	return out
}

func (s *OrderedMap[k, v]) ForEach(w typ.Tracker[k, v]) {
	e := s.uniques
	for _, ref := range s.order {
		key := *ref
		w(key, e[key])
	}
}


func (s *OrderedMap[k, v]) Len() int {
	return len(s.order)
}

func (s *OrderedMap[k, v]) Contains(key k) bool {
	_, ok := s.uniques[key]
	return ok
}


