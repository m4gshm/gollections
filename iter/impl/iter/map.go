package iter

import (
	"reflect"

	"github.com/m4gshm/container/K"
	"github.com/m4gshm/container/typ"
)

func NewKV[k comparable, v any](elements map[k]v) *KV[k, v] {
	refVal := reflect.ValueOf(elements)
	return &KV[k, v]{elements: elements, refVal: refVal, iter: refVal.MapRange()}
}

func NewOrderedKV[k comparable, v any](order []*k, uniques map[k]v) *OrderedKV[k, v] {
	return &OrderedKV[k, v]{elements: New(order), uniques: uniques}
}

type KV[k comparable, v any] struct {
	elements map[k]v
	iter     *reflect.MapIter
	refVal   reflect.Value
	err      error
}

var _ typ.Iterator[*typ.KV[interface{}, interface{}]] = (*KV[interface{}, interface{}])(nil)

func (s *KV[k, v]) HasNext() bool {
	next := s.iter.Next()
	if !next {
		s.err = Exhausted
	}
	return next
}

func (s *KV[k, v]) Get() *typ.KV[k, v] {
	if err := s.err; err != nil {
		panic(err)
	}
	key := s.iter.Key().Interface().(k)
	return K.V(key, s.elements[key])
}

func (s *KV[k, v]) Err() error {
	return s.err
}

func (s *KV[k, v]) Reset() {
	s.iter.Reset(s.refVal)
	s.err = nil
}

type OrderedKV[k comparable, v any] struct {
	elements *Iter[*k]
	uniques  map[k]v
}

var _ typ.Iterator[*typ.KV[any, any]] = (*OrderedKV[any, any])(nil)

func (s *OrderedKV[k, v]) HasNext() bool {
	return s.elements.HasNext()
}

func (s *OrderedKV[k, v]) Get() *typ.KV[k, v] {
	key := *s.elements.Get()
	return K.V(key, s.uniques[key])
}

func (s *OrderedKV[k, v]) Err() error {
	return s.elements.Err()
}
