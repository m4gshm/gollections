package immutable

import (
	"github.com/m4gshm/container/K"
	"github.com/m4gshm/container/iter/impl/iter"
	"github.com/m4gshm/container/typ"
)

func NewMap[k comparable, v any](values ...typ.KV[k, v]) *OrderedMap[k, v] {
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

var _ typ.Walk[*typ.KV[int, any]] = (*OrderedMap[int, any])(nil)
var _ typ.Container[*typ.KV[interface{}, interface{}], int] = (*OrderedMap[interface{}, interface{}])(nil)
// var _ fmt.Stringer = (*OrderedMap[interface{}, interface{}])(nil)

func (s *OrderedMap[k, v]) Begin() typ.Iterator[*typ.KV[k, v]] {
	return iter.NewOrderKV(s.order, s.uniques)
}

func (s *OrderedMap[k, v]) Values() []*typ.KV[k, v] {
	out := make([]*typ.KV[k, v], len(s.order))
	for i, keyRef := range s.order {
		key := *keyRef
		out[i] = K.V(key, s.uniques[key])
	}
	return out
}

func (s *OrderedMap[k, v]) Len() int {
	return len(s.order)
}
