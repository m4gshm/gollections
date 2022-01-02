package iter

import (
	"reflect"

	"github.com/m4gshm/container/K"
	"github.com/m4gshm/container/typ"
)

func NewKV[k comparable, v any](elements map[k]v) *KVIter[k, v] {
	refVal := reflect.ValueOf(elements)
	return &KVIter[k, v]{elements: elements, refVal: refVal, iter: refVal.MapRange()}
}

func NewOrderKV[k comparable, v any](order []*k, uniques map[k]v) *OrderedKVIter[k, v] {
	return &OrderedKVIter[k, v]{order: New(order), uniques: uniques}
}

type KVIter[k comparable, v any] struct {
	elements map[k]v
	iter     *reflect.MapIter
	refVal   reflect.Value
	err      error
}

var _ typ.Iterator[*typ.KV[interface{}, interface{}]] = (*KVIter[interface{}, interface{}])(nil)

func (s *KVIter[k, v]) HasNext() bool {
	next := s.iter.Next()
	if !next {
		s.err = Exhausted
	}
	return next
}

func (s *KVIter[k, v]) Get() *typ.KV[k, v] {
	if err := s.err; err != nil {
		panic(err)
	}
	key := s.iter.Key().Interface().(k)
	return K.V(key, s.elements[key])
}

func (s *KVIter[k, v]) Err() error {
	return s.err
}

func (s *KVIter[k, v]) Reset() {
	s.iter.Reset(s.refVal)
	s.err = nil
}

type OrderedKVIter[k comparable, v any] struct {
	order   *Iter[*k]
	uniques map[k]v
}

var _ typ.Iterator[*typ.KV[any, any]] = (*OrderedKVIter[any, any])(nil)

func (s *OrderedKVIter[k, v]) HasNext() bool {
	return s.order.HasNext()
}

func (s *OrderedKVIter[k, v]) Get() *typ.KV[k, v] {
	key := *s.order.Get()
	return K.V(key, s.uniques[key])
}

func (s *OrderedKVIter[k, v]) Err() error {
	return s.order.Err()
}
